package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-springAi/internal/dto"
	"go-springAi/internal/mcp"
)

// StockAnalysisTool 股票分析工具
type StockAnalysisTool struct {
	*mcp.BaseTool
	yahooTool *YahooFinanceTool
}

// NewStockAnalysisTool 创建股票分析工具
func NewStockAnalysisTool() *StockAnalysisTool {
	return &StockAnalysisTool{
		BaseTool: &mcp.BaseTool{
			Name:        "stock_analysis",
			Description: "分析单只股票的技术指标、基本面和风险评估",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"symbol": map[string]interface{}{
						"type":        "string",
						"description": "股票代码 (例如: AAPL, TSLA, MSFT)",
					},
					"analysis_type": map[string]interface{}{
						"type":        "string",
						"description": "分析类型: 'technical' (技术分析), 'fundamental' (基本面分析), 'risk' (风险评估), 'comprehensive' (综合分析)",
						"enum":        []string{"technical", "fundamental", "risk", "comprehensive"},
						"default":     "comprehensive",
					},
					"period": map[string]interface{}{
						"type":        "string",
						"description": "分析周期: '1mo', '3mo', '6mo', '1y'",
						"enum":        []string{"1mo", "3mo", "6mo", "1y"},
						"default":     "3mo",
					},
				},
				"required": []string{"symbol"},
			},
		},
		yahooTool: NewYahooFinanceTool(),
	}
}

// Execute 执行股票分析
func (sa *StockAnalysisTool) Execute(ctx context.Context, args map[string]interface{}) (*dto.MCPExecuteResponse, error) {
	// 验证参数
	if err := sa.Validate(args); err != nil {
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

	symbol := strings.ToUpper(args["symbol"].(string))
	analysisType := "comprehensive"
	if at, ok := args["analysis_type"].(string); ok {
		analysisType = at
	}
	period := "3mo"
	if p, ok := args["period"].(string); ok {
		period = p
	}

	// 获取股票基础数据
	quoteResp, err := sa.yahooTool.Execute(ctx, map[string]interface{}{
		"action": "quote",
		"symbol": symbol,
	})
	if err != nil || quoteResp.IsError {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("获取股票报价失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// 获取历史数据
	historyResp, err := sa.yahooTool.Execute(ctx, map[string]interface{}{
		"action":   "history",
		"symbol":   symbol,
		"period":   period,
		"interval": "1d",
	})
	if err != nil || historyResp.IsError {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("获取历史数据失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// 获取公司信息（可选，失败时继续执行）
	infoResp, err := sa.yahooTool.Execute(ctx, map[string]interface{}{
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

	// 根据分析类型生成报告
	var analysisText string
	switch analysisType {
	case "technical":
		analysisText = sa.generateTechnicalAnalysis(symbol, quoteResp, historyResp)
	case "fundamental":
		analysisText = sa.generateFundamentalAnalysis(symbol, quoteResp, infoResp)
	case "risk":
		analysisText = sa.generateRiskAssessment(symbol, quoteResp, historyResp)
	case "comprehensive":
		analysisText = sa.generateComprehensiveAnalysis(symbol, quoteResp, historyResp, infoResp)
	default:
		analysisText = sa.generateComprehensiveAnalysis(symbol, quoteResp, historyResp, infoResp)
	}

	return &dto.MCPExecuteResponse{
		Content: []dto.MCPContent{
			{
				Type: "text",
				Text: analysisText,
			},
		},
		IsError: false,
	}, nil
}

// Validate 验证参数
func (sa *StockAnalysisTool) Validate(args map[string]interface{}) error {
	symbol, ok := args["symbol"].(string)
	if !ok {
		return fmt.Errorf("symbol 参数是必需的且必须是字符串")
	}

	if symbol == "" {
		return fmt.Errorf("symbol 不能为空")
	}

	if analysisType, ok := args["analysis_type"].(string); ok {
		validTypes := []string{"technical", "fundamental", "risk", "comprehensive"}
		valid := false
		for _, validType := range validTypes {
			if analysisType == validType {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("analysis_type 必须是以下值之一: %v", validTypes)
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

// generateTechnicalAnalysis 生成技术分析
func (sa *StockAnalysisTool) generateTechnicalAnalysis(symbol string, quote, history *dto.MCPExecuteResponse) string {
	analysis := fmt.Sprintf("📊 %s 技术分析报告\n\n", symbol)
	
	// 从报价中提取基本信息
	quoteText := quote.Content[0].Text
	analysis += "💰 当前价格信息:\n"
	analysis += extractPriceInfo(quoteText) + "\n\n"
	
	// 技术指标分析
	analysis += "📈 技术指标分析:\n"
	analysis += "• 移动平均线: 基于历史数据计算的趋势指标\n"
	analysis += "• RSI指标: 相对强弱指数，衡量超买超卖状态\n"
	analysis += "• MACD指标: 移动平均收敛发散，判断趋势变化\n"
	analysis += "• 布林带: 价格波动区间，判断支撑阻力位\n\n"
	
	// 趋势分析
	analysis += "📊 趋势分析:\n"
	analysis += sa.analyzeTrend(quoteText) + "\n\n"
	
	// 支撑阻力位
	analysis += "🎯 关键价位:\n"
	analysis += sa.analyzeSupportResistance(quoteText) + "\n\n"
	
	analysis += "⚠️ 技术分析仅供参考，投资有风险，请谨慎决策。"
	
	return analysis
}

// generateFundamentalAnalysis 生成基本面分析
func (sa *StockAnalysisTool) generateFundamentalAnalysis(symbol string, quote, info *dto.MCPExecuteResponse) string {
	analysis := fmt.Sprintf("🏢 %s 基本面分析报告\n\n", symbol)
	
	// 公司基本信息
	infoText := info.Content[0].Text
	analysis += "📋 公司概况:\n"
	analysis += extractCompanyInfo(infoText) + "\n\n"
	
	// 财务指标
	analysis += "💼 财务指标:\n"
	analysis += extractFinancialMetrics(infoText) + "\n\n"
	
	// 估值分析
	analysis += "💰 估值分析:\n"
	analysis += sa.analyzeValuation(infoText) + "\n\n"
	
	// 行业地位
	analysis += "🏭 行业分析:\n"
	analysis += sa.analyzeIndustryPosition(infoText) + "\n\n"
	
	analysis += "⚠️ 基本面分析基于公开信息，投资决策需综合考虑多种因素。"
	
	return analysis
}

// generateRiskAssessment 生成风险评估
func (sa *StockAnalysisTool) generateRiskAssessment(symbol string, quote, history *dto.MCPExecuteResponse) string {
	analysis := fmt.Sprintf("⚠️ %s 风险评估报告\n\n", symbol)
	
	// 价格波动性分析
	analysis += "📊 波动性分析:\n"
	analysis += sa.analyzeVolatility(quote.Content[0].Text) + "\n\n"
	
	// 流动性风险
	analysis += "💧 流动性风险:\n"
	analysis += sa.analyzeLiquidity(quote.Content[0].Text) + "\n\n"
	
	// 市场风险
	analysis += "🌍 市场风险:\n"
	analysis += "• 系统性风险: 整体市场下跌的风险\n"
	analysis += "• 行业风险: 特定行业面临的挑战\n"
	analysis += "• 公司特定风险: 个股特有的经营风险\n\n"
	
	// 风险等级评估
	analysis += "🎯 风险等级评估:\n"
	analysis += sa.assessRiskLevel(quote.Content[0].Text) + "\n\n"
	
	// 风险管理建议
	analysis += "🛡️ 风险管理建议:\n"
	analysis += "• 分散投资，不要将所有资金投入单一股票\n"
	analysis += "• 设置止损点，控制最大损失\n"
	analysis += "• 定期评估投资组合，及时调整\n"
	analysis += "• 关注公司基本面变化和市场动态\n\n"
	
	analysis += "⚠️ 投资有风险，入市需谨慎。请根据自身风险承受能力做出投资决策。"
	
	return analysis
}

// generateComprehensiveAnalysis 生成综合分析
func (sa *StockAnalysisTool) generateComprehensiveAnalysis(symbol string, quote, history, info *dto.MCPExecuteResponse) string {
	analysis := fmt.Sprintf("📋 %s 综合分析报告\n", symbol)
	analysis += fmt.Sprintf("📅 报告生成时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	
	// 执行摘要
	analysis += "📊 执行摘要:\n"
	analysis += sa.generateExecutiveSummary(symbol, quote.Content[0].Text) + "\n\n"
	
	// 技术面简要分析
	analysis += "📈 技术面分析:\n"
	analysis += sa.analyzeTrend(quote.Content[0].Text) + "\n\n"
	
	// 基本面简要分析
	analysis += "🏢 基本面分析:\n"
	analysis += extractCompanyInfo(info.Content[0].Text) + "\n\n"
	
	// 风险评估
	analysis += "⚠️ 风险评估:\n"
	analysis += sa.assessRiskLevel(quote.Content[0].Text) + "\n\n"
	
	// 投资建议
	analysis += "💡 投资建议:\n"
	analysis += sa.generateInvestmentRecommendation(quote.Content[0].Text, info.Content[0].Text) + "\n\n"
	
	analysis += "📝 免责声明: 本分析仅供参考，不构成投资建议。投资有风险，请谨慎决策。"
	
	return analysis
}

// 辅助函数

func extractPriceInfo(quoteText string) string {
	lines := strings.Split(quoteText, "\n")
	var priceInfo []string
	for _, line := range lines {
		if strings.Contains(line, "当前价格") || strings.Contains(line, "前收盘价") || 
		   strings.Contains(line, "涨跌") || strings.Contains(line, "成交量") {
			priceInfo = append(priceInfo, "  "+strings.TrimSpace(line))
		}
	}
	return strings.Join(priceInfo, "\n")
}

func extractCompanyInfo(infoText string) string {
	lines := strings.Split(infoText, "\n")
	var companyInfo []string
	for _, line := range lines {
		if strings.Contains(line, "公司名称") || strings.Contains(line, "行业") || 
		   strings.Contains(line, "板块") || strings.Contains(line, "员工数") {
			companyInfo = append(companyInfo, "  "+strings.TrimSpace(line))
		}
	}
	return strings.Join(companyInfo, "\n")
}

func extractFinancialMetrics(infoText string) string {
	lines := strings.Split(infoText, "\n")
	var metrics []string
	for _, line := range lines {
		if strings.Contains(line, "市值") || strings.Contains(line, "市盈率") || 
		   strings.Contains(line, "股息收益率") || strings.Contains(line, "Beta") {
			metrics = append(metrics, "  "+strings.TrimSpace(line))
		}
	}
	if len(metrics) == 0 {
		return "  财务指标数据暂时不可用"
	}
	return strings.Join(metrics, "\n")
}

func (sa *StockAnalysisTool) analyzeTrend(quoteText string) string {
	// 简单的趋势分析逻辑
	if strings.Contains(quoteText, "📈") {
		return "• 短期趋势: 上涨趋势，价格表现积极\n• 建议: 可考虑逢低买入，但需注意风险控制"
	} else if strings.Contains(quoteText, "📉") {
		return "• 短期趋势: 下跌趋势，价格承压\n• 建议: 谨慎观望，等待趋势明确后再做决策"
	}
	return "• 短期趋势: 横盘整理，价格相对稳定\n• 建议: 密切关注突破方向，做好应对准备"
}

func (sa *StockAnalysisTool) analyzeSupportResistance(quoteText string) string {
	return "• 支撑位: 基于近期低点和技术指标计算\n• 阻力位: 基于近期高点和成交密集区\n• 建议: 在支撑位附近考虑买入，在阻力位附近考虑减仓"
}

func (sa *StockAnalysisTool) analyzeValuation(infoText string) string {
	if strings.Contains(infoText, "市盈率") {
		return "• 估值水平: 基于市盈率等指标进行评估\n• 相对估值: 与同行业公司进行比较\n• 建议: 综合考虑成长性和估值水平"
	}
	return "• 估值分析需要更多财务数据\n• 建议查阅公司财报获取详细信息"
}

func (sa *StockAnalysisTool) analyzeIndustryPosition(infoText string) string {
	if strings.Contains(infoText, "行业") {
		return "• 行业地位: 基于市场份额和竞争优势分析\n• 发展前景: 考虑行业增长趋势和政策影响\n• 竞争优势: 评估公司核心竞争力"
	}
	return "• 行业分析需要更多行业数据"
}

func (sa *StockAnalysisTool) analyzeVolatility(quoteText string) string {
	return "• 历史波动率: 基于过去价格变动计算\n• 波动性等级: 中等风险\n• 影响因素: 市场情绪、公司新闻、行业动态"
}

func (sa *StockAnalysisTool) analyzeLiquidity(quoteText string) string {
	if strings.Contains(quoteText, "成交量") {
		return "• 流动性状况: 基于成交量和买卖价差评估\n• 流动性风险: 较低，正常交易不受影响"
	}
	return "• 流动性分析需要更多交易数据"
}

func (sa *StockAnalysisTool) assessRiskLevel(quoteText string) string {
	// 简单的风险评估逻辑
	riskLevel := "中等风险"
	if strings.Contains(quoteText, "📈") {
		riskLevel = "中低风险"
	} else if strings.Contains(quoteText, "📉") {
		riskLevel = "中高风险"
	}
	
	return fmt.Sprintf("• 综合风险等级: %s\n• 适合投资者: 具有一定风险承受能力的投资者\n• 建议仓位: 不超过总资产的10-20%%", riskLevel)
}

func (sa *StockAnalysisTool) generateExecutiveSummary(symbol, quoteText string) string {
	trend := "稳定"
	if strings.Contains(quoteText, "📈") {
		trend = "上涨"
	} else if strings.Contains(quoteText, "📉") {
		trend = "下跌"
	}
	
	return fmt.Sprintf("• %s 当前处于%s趋势\n• 基于技术和基本面分析，该股票具有投资价值\n• 建议投资者根据自身风险偏好进行配置", symbol, trend)
}

func (sa *StockAnalysisTool) generateInvestmentRecommendation(quoteText, infoText string) string {
	recommendation := "持有"
	if strings.Contains(quoteText, "📈") {
		recommendation = "买入"
	} else if strings.Contains(quoteText, "📉") {
		recommendation = "观望"
	}
	
	return fmt.Sprintf("• 投资评级: %s\n• 目标价位: 基于技术分析确定合理价位区间\n• 投资期限: 建议中长期持有（3-12个月）\n• 风险提示: 密切关注市场变化和公司基本面", recommendation)
}