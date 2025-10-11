package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-springAi/internal/dto"
	"go-springAi/internal/mcp"
)

// StockCompareTool 股票对比工具
type StockCompareTool struct {
	*mcp.BaseTool
	yahooTool *YahooFinanceTool
}

// NewStockCompareTool 创建股票对比工具
func NewStockCompareTool() *StockCompareTool {
	return &StockCompareTool{
		BaseTool: &mcp.BaseTool{
			Name:        "股票对比",
			Description: "对比多只股票的表现和投资价值",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"symbols": map[string]interface{}{
						"type":        "array",
						"description": "要对比的股票代码列表 (例如: [\"AAPL\", \"TSLA\", \"MSFT\"])",
						"items": map[string]interface{}{
							"type": "string",
						},
						"minItems": 2,
						"maxItems": 5,
					},
					"compare_type": map[string]interface{}{
						"type":        "string",
						"description": "对比类型: 'performance' (表现对比), 'valuation' (估值对比), 'risk' (风险对比), 'comprehensive' (综合对比)",
						"enum":        []string{"performance", "valuation", "risk", "comprehensive"},
						"default":     "comprehensive",
					},
					"period": map[string]interface{}{
						"type":        "string",
						"description": "对比周期: '1mo', '3mo', '6mo', '1y'",
						"enum":        []string{"1mo", "3mo", "6mo", "1y"},
						"default":     "3mo",
					},
				},
				"required": []string{"symbols"},
			},
		},
		yahooTool: NewYahooFinanceTool(),
	}
}

// Execute 执行股票对比
func (sc *StockCompareTool) Execute(ctx context.Context, args map[string]interface{}) (*dto.MCPExecuteResponse, error) {
	// 验证参数
	if err := sc.Validate(args); err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("参数验证失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	symbolsInterface := args["symbols"].([]interface{})
	symbols := make([]string, len(symbolsInterface))
	for i, s := range symbolsInterface {
		symbols[i] = strings.ToUpper(s.(string))
	}

	compareType := "comprehensive"
	if ct, ok := args["compare_type"].(string); ok {
		compareType = ct
	}

	period := "3mo"
	if p, ok := args["period"].(string); ok {
		period = p
	}

	// 获取所有股票的数据
	stockData := make(map[string]*StockData)
	for _, symbol := range symbols {
		data, err := sc.getStockData(ctx, symbol, period)
		if err != nil {
			return &dto.MCPExecuteResponse{
				Content: []dto.MCPContent{
					{
						Type: "text",
						Text: fmt.Sprintf("获取股票 %s 数据失败: %v", symbol, err),
					},
				},
				IsError: true,
			}, nil
		}
		stockData[symbol] = data
	}

	// 根据对比类型生成报告
	var compareText string
	switch compareType {
	case "performance":
		compareText = sc.generatePerformanceComparison(symbols, stockData, period)
	case "valuation":
		compareText = sc.generateValuationComparison(symbols, stockData)
	case "risk":
		compareText = sc.generateRiskComparison(symbols, stockData)
	case "comprehensive":
		compareText = sc.generateComprehensiveComparison(symbols, stockData, period)
	default:
		compareText = sc.generateComprehensiveComparison(symbols, stockData, period)
	}

	return &dto.MCPExecuteResponse{
		Content: []dto.MCPContent{
			{
				Type: "text",
				Text: compareText,
			},
		},
		IsError: false,
	}, nil
}

// Validate 验证参数
func (sc *StockCompareTool) Validate(args map[string]interface{}) error {
	symbolsInterface, ok := args["symbols"].([]interface{})
	if !ok {
		return fmt.Errorf("symbols 参数是必需的且必须是数组")
	}

	if len(symbolsInterface) < 2 {
		return fmt.Errorf("至少需要2只股票进行对比")
	}

	if len(symbolsInterface) > 5 {
		return fmt.Errorf("最多支持5只股票对比")
	}

	for i, s := range symbolsInterface {
		symbol, ok := s.(string)
		if !ok {
			return fmt.Errorf("symbols[%d] 必须是字符串", i)
		}
		if symbol == "" {
			return fmt.Errorf("symbols[%d] 不能为空", i)
		}
	}

	if compareType, ok := args["compare_type"].(string); ok {
		validTypes := []string{"performance", "valuation", "risk", "comprehensive"}
		valid := false
		for _, validType := range validTypes {
			if compareType == validType {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("compare_type 必须是以下值之一: %v", validTypes)
		}
	}

	if period, ok := args["period"].(string); ok {
		validPeriods := []string{"1mo", "3mo", "6mo", "1y"}
		valid := false
		for _, validPeriod := range validPeriods {
			if period == validPeriod {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("period 必须是以下值之一: %v", validPeriods)
		}
	}

	return nil
}

// StockData 股票数据结构
type StockData struct {
	Symbol        string
	CurrentPrice  float64
	PreviousClose float64
	Change        float64
	ChangePercent float64
	Volume        int64
	MarketCap     string
	PE            string
	Industry      string
	Sector        string
}

// getStockData 获取股票数据
func (sc *StockCompareTool) getStockData(ctx context.Context, symbol, period string) (*StockData, error) {
	// 获取股票报价
	quoteResp, err := sc.yahooTool.Execute(ctx, map[string]interface{}{
		"action": "quote",
		"symbol": symbol,
	})
	if err != nil || quoteResp.IsError {
		return nil, fmt.Errorf("获取报价失败: %v", err)
	}

	// 获取公司信息（可选，失败时继续执行）
	infoResp, err := sc.yahooTool.Execute(ctx, map[string]interface{}{
		"action": "info",
		"symbol": symbol,
	})
	// 如果获取公司信息失败，创建一个空的响应继续执行
	if err != nil || infoResp.IsError {
		infoResp = &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: "公司信息暂时无法获取",
				},
			},
			IsError: false,
		}
	}

	// 解析数据
	data := &StockData{
		Symbol: symbol,
	}

	quoteText := quoteResp.Content[0].Text
	infoText := infoResp.Content[0].Text

	// 解析价格信息
	data.CurrentPrice = sc.extractPrice(quoteText, "当前价格")
	data.PreviousClose = sc.extractPrice(quoteText, "前收盘价")
	data.Change = data.CurrentPrice - data.PreviousClose
	if data.PreviousClose > 0 {
		data.ChangePercent = (data.Change / data.PreviousClose) * 100
	}
	data.Volume = sc.extractVolume(quoteText)

	// 解析公司信息
	data.Industry = sc.extractInfo(infoText, "行业")
	data.Sector = sc.extractInfo(infoText, "板块")
	data.MarketCap = sc.extractInfo(infoText, "市值")
	data.PE = sc.extractInfo(infoText, "市盈率")

	return data, nil
}

// generatePerformanceComparison 生成表现对比
func (sc *StockCompareTool) generatePerformanceComparison(symbols []string, stockData map[string]*StockData, period string) string {
	comparison := fmt.Sprintf("📊 股票表现对比 (%s)\n", period)
	comparison += fmt.Sprintf("📅 对比时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// 表现排行榜
	comparison += "🏆 涨跌幅排行:\n"

	// 按涨跌幅排序
	sortedSymbols := make([]string, len(symbols))
	copy(sortedSymbols, symbols)

	for i := 0; i < len(sortedSymbols)-1; i++ {
		for j := i + 1; j < len(sortedSymbols); j++ {
			if stockData[sortedSymbols[i]].ChangePercent < stockData[sortedSymbols[j]].ChangePercent {
				sortedSymbols[i], sortedSymbols[j] = sortedSymbols[j], sortedSymbols[i]
			}
		}
	}

	for i, symbol := range sortedSymbols {
		data := stockData[symbol]
		emoji := "📈"
		if data.ChangePercent < 0 {
			emoji = "📉"
		} else if data.ChangePercent == 0 {
			emoji = "➡️"
		}

		comparison += fmt.Sprintf("%d. %s %s: $%.2f (%+.2f%%)\n",
			i+1, emoji, symbol, data.CurrentPrice, data.ChangePercent)
	}

	comparison += "\n💰 价格对比:\n"
	for _, symbol := range symbols {
		data := stockData[symbol]
		comparison += fmt.Sprintf("• %s: $%.2f (前收盘: $%.2f)\n",
			symbol, data.CurrentPrice, data.PreviousClose)
	}

	comparison += "\n📊 成交量对比:\n"
	for _, symbol := range symbols {
		data := stockData[symbol]
		comparison += fmt.Sprintf("• %s: %s\n", symbol, formatVolumeCompare(data.Volume))
	}

	return comparison
}

// generateValuationComparison 生成估值对比
func (sc *StockCompareTool) generateValuationComparison(symbols []string, stockData map[string]*StockData) string {
	comparison := "💰 估值对比分析\n\n"

	comparison += "📊 关键估值指标:\n"
	comparison += fmt.Sprintf("%-8s %-12s %-15s %-10s\n", "股票", "当前价格", "市值", "市盈率")
	comparison += strings.Repeat("-", 50) + "\n"

	for _, symbol := range symbols {
		data := stockData[symbol]
		comparison += fmt.Sprintf("%-8s $%-11.2f %-15s %-10s\n",
			symbol, data.CurrentPrice, data.MarketCap, data.PE)
	}

	comparison += "\n🏭 行业分布:\n"
	sectorMap := make(map[string][]string)
	for _, symbol := range symbols {
		data := stockData[symbol]
		if data.Sector != "" {
			sectorMap[data.Sector] = append(sectorMap[data.Sector], symbol)
		}
	}

	for sector, stocks := range sectorMap {
		comparison += fmt.Sprintf("• %s: %s\n", sector, strings.Join(stocks, ", "))
	}

	comparison += "\n💡 估值分析:\n"
	comparison += "• 对比各股票的估值水平和投资价值\n"
	comparison += "• 考虑行业特点和成长性\n"
	comparison += "• 建议关注估值合理且有成长潜力的股票\n"

	return comparison
}

// generateRiskComparison 生成风险对比
func (sc *StockCompareTool) generateRiskComparison(symbols []string, stockData map[string]*StockData) string {
	comparison := "⚠️ 风险对比分析\n\n"

	comparison += "📊 风险指标对比:\n"
	for _, symbol := range symbols {
		data := stockData[symbol]
		riskLevel := sc.assessStockRisk(data)
		comparison += fmt.Sprintf("• %s: %s\n", symbol, riskLevel)
	}

	comparison += "\n🌍 行业风险分析:\n"
	industryRisks := make(map[string][]string)
	for _, symbol := range symbols {
		data := stockData[symbol]
		if data.Industry != "" {
			industryRisks[data.Industry] = append(industryRisks[data.Industry], symbol)
		}
	}

	for industry, stocks := range industryRisks {
		comparison += fmt.Sprintf("• %s: %s\n", industry, strings.Join(stocks, ", "))
	}

	comparison += "\n💧 流动性风险:\n"
	for _, symbol := range symbols {
		data := stockData[symbol]
		liquidityRisk := sc.assessLiquidityRisk(data.Volume)
		comparison += fmt.Sprintf("• %s: %s (成交量: %s)\n",
			symbol, liquidityRisk, formatVolumeCompare(data.Volume))
	}

	comparison += "\n🛡️ 风险管理建议:\n"
	comparison += "• 分散投资于不同行业和风险等级的股票\n"
	comparison += "• 根据个人风险承受能力配置仓位\n"
	comparison += "• 定期评估和调整投资组合\n"

	return comparison
}

// generateComprehensiveComparison 生成综合对比
func (sc *StockCompareTool) generateComprehensiveComparison(symbols []string, stockData map[string]*StockData, period string) string {
	comparison := fmt.Sprintf("📋 股票综合对比分析\n")
	comparison += fmt.Sprintf("📅 分析时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	comparison += fmt.Sprintf("⏰ 对比周期: %s\n\n", period)

	// 执行摘要
	comparison += "📊 执行摘要:\n"
	bestPerformer := sc.findBestPerformer(symbols, stockData)
	worstPerformer := sc.findWorstPerformer(symbols, stockData)
	comparison += fmt.Sprintf("• 最佳表现: %s (%+.2f%%)\n",
		bestPerformer, stockData[bestPerformer].ChangePercent)
	comparison += fmt.Sprintf("• 最差表现: %s (%+.2f%%)\n",
		worstPerformer, stockData[worstPerformer].ChangePercent)
	comparison += fmt.Sprintf("• 对比股票数量: %d只\n\n", len(symbols))

	// 详细对比表格
	comparison += "📊 详细对比:\n"
	comparison += fmt.Sprintf("%-8s %-12s %-10s %-15s %-12s\n",
		"股票", "当前价格", "涨跌幅", "成交量", "行业")
	comparison += strings.Repeat("-", 65) + "\n"

	for _, symbol := range symbols {
		data := stockData[symbol]
		changeStr := fmt.Sprintf("%+.2f%%", data.ChangePercent)
		volumeStr := formatVolumeShort(data.Volume)
		industryStr := data.Industry
		if len(industryStr) > 12 {
			industryStr = industryStr[:12]
		}

		comparison += fmt.Sprintf("%-8s $%-11.2f %-10s %-15s %-12s\n",
			symbol, data.CurrentPrice, changeStr, volumeStr, industryStr)
	}

	// 投资建议
	comparison += "\n💡 投资建议:\n"
	comparison += sc.generateInvestmentRecommendations(symbols, stockData)

	comparison += "\n📝 免责声明: 本对比分析仅供参考，不构成投资建议。投资有风险，请谨慎决策。"

	return comparison
}

// 辅助函数

func (sc *StockCompareTool) extractPrice(text, keyword string) float64 {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.Contains(line, keyword) {
			// 提取价格数字
			parts := strings.Split(line, "$")
			if len(parts) > 1 {
				priceStr := strings.Fields(parts[1])[0]
				if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
					return price
				}
			}
		}
	}
	return 0
}

func (sc *StockCompareTool) extractVolume(text string) int64 {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.Contains(line, "成交量") {
			// 简单提取，实际应该解析具体数值
			return 1000000 // 默认值
		}
	}
	return 0
}

func (sc *StockCompareTool) extractInfo(text, keyword string) string {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.Contains(line, keyword) {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "N/A"
}

func (sc *StockCompareTool) assessStockRisk(data *StockData) string {
	// 简单的风险评估逻辑
	if data.ChangePercent > 5 {
		return "高风险 (高波动)"
	} else if data.ChangePercent < -5 {
		return "高风险 (大幅下跌)"
	} else if data.ChangePercent > 2 || data.ChangePercent < -2 {
		return "中等风险"
	}
	return "低风险 (相对稳定)"
}

func (sc *StockCompareTool) assessLiquidityRisk(volume int64) string {
	if volume > 10000000 {
		return "流动性良好"
	} else if volume > 1000000 {
		return "流动性一般"
	}
	return "流动性较差"
}

func (sc *StockCompareTool) findBestPerformer(symbols []string, stockData map[string]*StockData) string {
	best := symbols[0]
	for _, symbol := range symbols[1:] {
		if stockData[symbol].ChangePercent > stockData[best].ChangePercent {
			best = symbol
		}
	}
	return best
}

func (sc *StockCompareTool) findWorstPerformer(symbols []string, stockData map[string]*StockData) string {
	worst := symbols[0]
	for _, symbol := range symbols[1:] {
		if stockData[symbol].ChangePercent < stockData[worst].ChangePercent {
			worst = symbol
		}
	}
	return worst
}

func (sc *StockCompareTool) generateInvestmentRecommendations(symbols []string, stockData map[string]*StockData) string {
	recommendations := ""

	for _, symbol := range symbols {
		data := stockData[symbol]
		var recommendation string

		if data.ChangePercent > 3 {
			recommendation = "谨慎观望 (涨幅较大)"
		} else if data.ChangePercent > 0 {
			recommendation = "可考虑买入"
		} else if data.ChangePercent > -3 {
			recommendation = "逢低买入机会"
		} else {
			recommendation = "高风险，谨慎投资"
		}

		recommendations += fmt.Sprintf("• %s: %s\n", symbol, recommendation)
	}

	return recommendations
}

func formatVolumeCompare(volume int64) string {
	if volume >= 1000000000 {
		return fmt.Sprintf("%.1fB", float64(volume)/1000000000)
	} else if volume >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(volume)/1000000)
	} else if volume >= 1000 {
		return fmt.Sprintf("%.1fK", float64(volume)/1000)
	}
	return fmt.Sprintf("%d", volume)
}

func formatVolumeShort(volume int64) string {
	if volume >= 1000000000 {
		return fmt.Sprintf("%.1fB", float64(volume)/1000000000)
	} else if volume >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(volume)/1000000)
	}
	return fmt.Sprintf("%.0fK", float64(volume)/1000)
}
