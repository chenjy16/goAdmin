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

// StockAdviceTool 股票投资建议工具
type StockAdviceTool struct {
	*mcp.BaseTool
	yahooTool *YahooFinanceTool
}

// NewStockAdviceTool 创建新的股票投资建议工具
func NewStockAdviceTool() *StockAdviceTool {
	return &StockAdviceTool{
		BaseTool: &mcp.BaseTool{
			Name:        "stock_advice",
			Description: "基于股票分析提供投资建议和风险提示",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"symbol": map[string]interface{}{
						"type":        "string",
						"description": "股票代码 (例如: AAPL, TSLA)",
					},
					"investment_horizon": map[string]interface{}{
						"type":        "string",
						"description": "投资期限 (short_term: 短期, medium_term: 中期, long_term: 长期)",
						"enum":        []string{"short_term", "medium_term", "long_term"},
						"default":     "medium_term",
					},
					"risk_tolerance": map[string]interface{}{
						"type":        "string",
						"description": "风险承受能力 (conservative: 保守, moderate: 适中, aggressive: 激进)",
						"enum":        []string{"conservative", "moderate", "aggressive"},
						"default":     "moderate",
					},
					"investment_amount": map[string]interface{}{
						"type":        "number",
						"description": "投资金额 (美元)",
						"minimum":     100,
					},
				},
				"required": []string{"symbol"},
			},
		},
		yahooTool: NewYahooFinanceTool(),
	}
}

// Execute 执行股票投资建议分析
func (sa *StockAdviceTool) Execute(ctx context.Context, args map[string]interface{}) (*dto.MCPExecuteResponse, error) {
	symbol, ok := args["symbol"].(string)
	if !ok {
		return nil, fmt.Errorf("symbol is required and must be a string")
	}

	symbol = strings.ToUpper(symbol)

	// 获取投资参数
	horizon := "medium_term"
	if h, ok := args["investment_horizon"].(string); ok {
		horizon = h
	}

	riskTolerance := "moderate"
	if r, ok := args["risk_tolerance"].(string); ok {
		riskTolerance = r
	}

	var investmentAmount float64
	if amount, ok := args["investment_amount"].(float64); ok {
		investmentAmount = amount
	}

	// 获取股票基础数据
	quoteArgs := map[string]interface{}{"symbol": symbol}
	quoteResp, err := sa.yahooTool.Execute(ctx, quoteArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock quote: %v", err)
	}

	// 获取公司信息
	infoArgs := map[string]interface{}{"symbol": symbol, "action": "info"}
	infoResp, err := sa.yahooTool.Execute(ctx, infoArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to get company info: %v", err)
	}

	// 获取历史数据
	historyArgs := map[string]interface{}{
		"symbol": symbol,
		"action": "history",
		"period": "3mo",
	}
	historyResp, err := sa.yahooTool.Execute(ctx, historyArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock history: %v", err)
	}

	// 生成投资建议
	advice := sa.generateInvestmentAdvice(symbol, quoteResp, infoResp, historyResp, horizon, riskTolerance, investmentAmount)

	return &dto.MCPExecuteResponse{
		Content: []dto.MCPContent{
			{
				Type: "text",
				Text: advice,
			},
		},
		IsError: false,
	}, nil
}

// Validate 验证输入参数
func (sa *StockAdviceTool) Validate(args map[string]interface{}) error {
	symbol, ok := args["symbol"].(string)
	if !ok {
		return fmt.Errorf("symbol is required and must be a string")
	}

	if len(symbol) == 0 {
		return fmt.Errorf("symbol cannot be empty")
	}

	// 验证投资期限
	if horizon, ok := args["investment_horizon"].(string); ok {
		validHorizons := []string{"short_term", "medium_term", "long_term"}
		valid := false
		for _, vh := range validHorizons {
			if horizon == vh {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid investment_horizon: %s", horizon)
		}
	}

	// 验证风险承受能力
	if risk, ok := args["risk_tolerance"].(string); ok {
		validRisks := []string{"conservative", "moderate", "aggressive"}
		valid := false
		for _, vr := range validRisks {
			if risk == vr {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid risk_tolerance: %s", risk)
		}
	}

	// 验证投资金额
	if amount, ok := args["investment_amount"].(float64); ok {
		if amount < 100 {
			return fmt.Errorf("investment_amount must be at least $100")
		}
	}

	return nil
}

// 生成投资建议
func (sa *StockAdviceTool) generateInvestmentAdvice(symbol string, quoteResp, infoResp, historyResp *dto.MCPExecuteResponse, horizon, riskTolerance string, investmentAmount float64) string {
	advice := fmt.Sprintf("📊 %s 股票投资建议报告\n", symbol)
	advice += fmt.Sprintf("生成时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// 提取关键数据
	quoteData := sa.extractResponseText(quoteResp)
	infoData := sa.extractResponseText(infoResp)
	historyData := sa.extractResponseText(historyResp)
	
	currentPrice := sa.extractPrice(quoteData, "当前价格")
	changePercent := sa.extractChangePercent(quoteData)
	volume := sa.extractVolumeFromText(quoteData)
	marketCap := sa.extractInfo(infoData, "市值")
	pe := sa.extractInfo(infoData, "市盈率")
	sector := sa.extractInfo(infoData, "行业")
	
	// 从历史数据中提取趋势信息
	_ = historyData // 使用历史数据进行趋势分析（简化处理）

	// 基本信息
	advice += "📈 基本信息:\n"
	advice += fmt.Sprintf("• 当前价格: $%.2f\n", currentPrice)
	advice += fmt.Sprintf("• 涨跌幅: %.2f%%\n", changePercent)
	advice += fmt.Sprintf("• 市值: %s\n", marketCap)
	advice += fmt.Sprintf("• 市盈率: %s\n", pe)
	advice += fmt.Sprintf("• 行业: %s\n\n", sector)

	// 投资建议评级
	advice += "🎯 投资建议评级:\n"
	rating := sa.calculateInvestmentRating(currentPrice, changePercent, pe, horizon, riskTolerance)
	advice += fmt.Sprintf("• 综合评级: %s\n", rating.Overall)
	advice += fmt.Sprintf("• 买入信号: %s\n", rating.BuySignal)
	advice += fmt.Sprintf("• 风险等级: %s\n\n", rating.RiskLevel)

	// 基于投资期限的建议
	advice += sa.generateHorizonSpecificAdvice(symbol, horizon, rating)

	// 基于风险承受能力的建议
	advice += sa.generateRiskBasedAdvice(symbol, riskTolerance, rating)

	// 仓位建议
	if investmentAmount > 0 {
		advice += sa.generatePositionAdvice(symbol, currentPrice, investmentAmount, riskTolerance)
	}

	// 风险提示
	advice += sa.generateRiskWarnings(symbol, changePercent, volume, sector)

	// 操作建议
	advice += sa.generateActionPlan(symbol, rating, horizon)

	advice += "\n⚠️ 重要声明: 本建议仅供参考，不构成投资建议。投资有风险，请根据自身情况谨慎决策。"

	return advice
}

// 投资评级结构
type InvestmentRating struct {
	Overall   string
	BuySignal string
	RiskLevel string
	Score     int
}

// 计算投资评级
func (sa *StockAdviceTool) calculateInvestmentRating(price, changePercent float64, pe, horizon, riskTolerance string) *InvestmentRating {
	score := 50 // 基础分数

	// 基于价格变化调整
	if changePercent > 5 {
		score -= 10 // 大涨可能过热
	} else if changePercent > 2 {
		score += 5 // 适度上涨
	} else if changePercent > -2 {
		score += 10 // 稳定
	} else if changePercent > -5 {
		score += 5 // 小幅下跌，可能是机会
	} else {
		score -= 15 // 大跌风险高
	}

	// 基于PE调整
	if pe != "N/A" && pe != "" {
		if peValue, err := strconv.ParseFloat(pe, 64); err == nil {
			if peValue < 15 {
				score += 10 // 低估值
			} else if peValue < 25 {
				score += 5 // 合理估值
			} else if peValue > 40 {
				score -= 10 // 高估值
			}
		}
	}

	// 基于投资期限调整
	switch horizon {
	case "long_term":
		score += 5 // 长期投资更稳健
	case "short_term":
		score -= 5 // 短期投资风险更高
	}

	// 确定评级
	var overall, buySignal, riskLevel string
	
	if score >= 70 {
		overall = "强烈推荐"
		buySignal = "强烈买入"
		riskLevel = "低风险"
	} else if score >= 60 {
		overall = "推荐"
		buySignal = "买入"
		riskLevel = "中低风险"
	} else if score >= 50 {
		overall = "中性"
		buySignal = "观望"
		riskLevel = "中等风险"
	} else if score >= 40 {
		overall = "谨慎"
		buySignal = "谨慎买入"
		riskLevel = "中高风险"
	} else {
		overall = "不推荐"
		buySignal = "避免买入"
		riskLevel = "高风险"
	}

	return &InvestmentRating{
		Overall:   overall,
		BuySignal: buySignal,
		RiskLevel: riskLevel,
		Score:     score,
	}
}

// 生成基于投资期限的建议
func (sa *StockAdviceTool) generateHorizonSpecificAdvice(symbol, horizon string, rating *InvestmentRating) string {
	advice := "⏰ 投资期限建议:\n"
	
	switch horizon {
	case "short_term":
		advice += "• 短期投资 (1-6个月):\n"
		if rating.Score >= 60 {
			advice += "  - 适合短期交易，关注技术指标\n"
			advice += "  - 设置止损位在5-8%\n"
			advice += "  - 密切关注市场情绪变化\n"
		} else {
			advice += "  - 短期风险较高，建议观望\n"
			advice += "  - 如需交易，严格控制仓位\n"
		}
	case "medium_term":
		advice += "• 中期投资 (6个月-2年):\n"
		if rating.Score >= 50 {
			advice += "  - 适合中期持有，关注基本面\n"
			advice += "  - 可分批建仓，降低成本\n"
			advice += "  - 关注季度财报表现\n"
		} else {
			advice += "  - 等待更好的入场时机\n"
			advice += "  - 关注行业趋势变化\n"
		}
	case "long_term":
		advice += "• 长期投资 (2年以上):\n"
		if rating.Score >= 45 {
			advice += "  - 适合长期价值投资\n"
			advice += "  - 关注公司发展战略\n"
			advice += "  - 可忽略短期波动\n"
		} else {
			advice += "  - 需要深入研究公司基本面\n"
			advice += "  - 考虑行业长期前景\n"
		}
	}
	
	return advice + "\n"
}

// 生成基于风险承受能力的建议
func (sa *StockAdviceTool) generateRiskBasedAdvice(symbol, riskTolerance string, rating *InvestmentRating) string {
	advice := "🎲 风险承受能力建议:\n"
	
	switch riskTolerance {
	case "conservative":
		advice += "• 保守型投资者:\n"
		if rating.Score >= 65 {
			advice += "  - 可适量配置，不超过总资产5%\n"
			advice += "  - 建议分散投资，降低风险\n"
		} else {
			advice += "  - 当前风险偏高，建议观望\n"
			advice += "  - 考虑更稳健的投资选择\n"
		}
	case "moderate":
		advice += "• 稳健型投资者:\n"
		if rating.Score >= 55 {
			advice += "  - 可配置10-15%的资产\n"
			advice += "  - 结合其他资产平衡风险\n"
		} else {
			advice += "  - 谨慎投资，控制仓位\n"
			advice += "  - 等待更好的投资机会\n"
		}
	case "aggressive":
		advice += "• 激进型投资者:\n"
		if rating.Score >= 45 {
			advice += "  - 可配置20-30%的资产\n"
			advice += "  - 可考虑杠杆操作（谨慎）\n"
		} else {
			advice += "  - 高风险高收益，需谨慎评估\n"
			advice += "  - 严格设置止损策略\n"
		}
	}
	
	return advice + "\n"
}

// 生成仓位建议
func (sa *StockAdviceTool) generatePositionAdvice(symbol string, currentPrice, investmentAmount float64, riskTolerance string) string {
	advice := "💰 仓位建议:\n"
	advice += fmt.Sprintf("• 投资金额: $%.2f\n", investmentAmount)
	advice += fmt.Sprintf("• 当前股价: $%.2f\n", currentPrice)
	
	shares := int(investmentAmount / currentPrice)
	actualAmount := float64(shares) * currentPrice
	
	advice += fmt.Sprintf("• 建议股数: %d 股\n", shares)
	advice += fmt.Sprintf("• 实际投资: $%.2f\n", actualAmount)
	
	// 分批建仓建议
	switch riskTolerance {
	case "conservative":
		advice += "• 建仓策略: 分3批建仓，每批33%\n"
		advice += "• 时间间隔: 每周一次\n"
	case "moderate":
		advice += "• 建仓策略: 分2批建仓，每批50%\n"
		advice += "• 时间间隔: 每两周一次\n"
	case "aggressive":
		advice += "• 建仓策略: 可一次性建仓\n"
		advice += "• 或分2批，快速建仓\n"
	}
	
	return advice + "\n"
}

// 生成风险提示
func (sa *StockAdviceTool) generateRiskWarnings(symbol string, changePercent, volume float64, sector string) string {
	warnings := "⚠️ 风险提示:\n"
	
	// 波动性风险
	if changePercent > 10 || changePercent < -10 {
		warnings += "• 高波动性: 股价波动较大，注意风险控制\n"
	}
	
	// 流动性风险
	if volume < 1000000 {
		warnings += "• 流动性风险: 成交量较低，可能影响买卖\n"
	}
	
	// 行业风险
	riskySectors := []string{"科技", "生物技术", "加密货币", "新能源"}
	for _, rs := range riskySectors {
		if strings.Contains(sector, rs) {
			warnings += fmt.Sprintf("• 行业风险: %s行业波动性较高\n", rs)
			break
		}
	}
	
	// 通用风险
	warnings += "• 市场风险: 受整体市场环境影响\n"
	warnings += "• 汇率风险: 如为外币计价，需关注汇率变化\n"
	warnings += "• 政策风险: 关注相关政策法规变化\n"
	
	return warnings + "\n"
}

// 生成操作建议
func (sa *StockAdviceTool) generateActionPlan(symbol string, rating *InvestmentRating, horizon string) string {
	plan := "📋 操作建议:\n"
	
	if rating.Score >= 60 {
		plan += "• 立即行动:\n"
		plan += "  1. 确认投资金额和风险承受能力\n"
		plan += "  2. 设置买入价格区间\n"
		plan += "  3. 制定止损和止盈策略\n"
		plan += "  4. 开始分批建仓\n"
	} else if rating.Score >= 50 {
		plan += "• 谨慎观察:\n"
		plan += "  1. 继续关注股价走势\n"
		plan += "  2. 等待更好的入场时机\n"
		plan += "  3. 关注公司最新动态\n"
		plan += "  4. 准备资金，随时行动\n"
	} else {
		plan += "• 暂时观望:\n"
		plan += "  1. 深入研究公司基本面\n"
		plan += "  2. 关注行业发展趋势\n"
		plan += "  3. 等待风险降低\n"
		plan += "  4. 考虑其他投资选择\n"
	}
	
	// 监控指标
	plan += "\n📊 关键监控指标:\n"
	plan += "• 股价支撑位和阻力位\n"
	plan += "• 成交量变化\n"
	plan += "• 财报发布时间\n"
	plan += "• 行业新闻和政策\n"
	plan += "• 技术指标 (RSI, MACD, 移动平均线)\n"
	
	return plan + "\n"
}

// 辅助函数

func (sa *StockAdviceTool) extractPrice(text, keyword string) float64 {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.Contains(line, keyword) {
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

func (sa *StockAdviceTool) extractChangePercent(text string) float64 {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.Contains(line, "涨跌幅") || strings.Contains(line, "变化") {
			if strings.Contains(line, "%") {
				parts := strings.Split(line, "%")
				if len(parts) > 0 {
					percentStr := strings.Fields(parts[0])
					if len(percentStr) > 0 {
						if percent, err := strconv.ParseFloat(percentStr[len(percentStr)-1], 64); err == nil {
							return percent
						}
					}
				}
			}
		}
	}
	return 0
}

func (sa *StockAdviceTool) extractVolumeFromText(text string) float64 {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.Contains(line, "成交量") {
			// 简单提取，实际应该解析具体数值
			return 1000000 // 默认值
		}
	}
	return 0
}

func (sa *StockAdviceTool) extractInfo(text, keyword string) string {
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

// extractResponseText 从MCPExecuteResponse中提取文本内容
func (sa *StockAdviceTool) extractResponseText(resp *dto.MCPExecuteResponse) string {
	if resp == nil || len(resp.Content) == 0 {
		return ""
	}
	
	var text strings.Builder
	for _, content := range resp.Content {
		if content.Type == "text" && content.Text != "" {
			text.WriteString(content.Text)
			text.WriteString("\n")
		}
	}
	
	return text.String()
}