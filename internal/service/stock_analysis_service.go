package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-springAi/internal/dto"
	"go-springAi/internal/mcp"

	"go.uber.org/zap"
)

// StockAnalysisService 股票分析服务
type StockAnalysisService struct {
	mcpClient mcp.InternalMCPClient
	logger    *zap.Logger
}

// NewStockAnalysisService 创建股票分析服务
func NewStockAnalysisService(mcpClient mcp.InternalMCPClient, logger *zap.Logger) *StockAnalysisService {
	service := &StockAnalysisService{
		mcpClient: mcpClient,
		logger:    logger,
	}
	
	// 自动初始化MCP客户端
	ctx := context.Background()
	initReq := &dto.MCPInitializeRequest{
		ProtocolVersion: "2024-11-05",
		Capabilities: dto.MCPCapabilities{
			Tools: &dto.MCPToolsCapability{
				ListChanged: true,
			},
		},
		ClientInfo: dto.MCPClientInfo{
			Name:    "stock-analysis-service",
			Version: "1.0.0",
		},
	}
	
	_, err := mcpClient.Initialize(ctx, initReq)
	if err != nil {
		logger.Error("Failed to initialize MCP client for stock analysis service", zap.Error(err))
	} else {
		logger.Info("MCP client initialized successfully for stock analysis service")
	}
	
	return service
}

// AnalyzeStock 分析单只股票
func (s *StockAnalysisService) AnalyzeStock(ctx context.Context, req *dto.StockAnalysisRequest) (*dto.StockAnalysisResponse, error) {
	s.logger.Info("开始分析股票", zap.String("symbol", req.Symbol), zap.String("analysis_type", req.AnalysisType))

	// 1. 获取股票基本信息
	quote, err := s.getStockQuote(ctx, req.Symbol)
	if err != nil {
		return nil, fmt.Errorf("获取股票报价失败: %w", err)
	}

	// 2. 获取历史数据
	period := req.Period
	if period == "" {
		period = "3mo" // 默认3个月
	}
	history, err := s.getStockHistory(ctx, req.Symbol, period, "1d")
	if err != nil {
		s.logger.Warn("获取历史数据失败", zap.Error(err))
	}

	// 3. 获取公司信息
	companyInfo, err := s.getStockInfo(ctx, req.Symbol)
	if err != nil {
		s.logger.Warn("获取公司信息失败", zap.Error(err))
	}

	// 4. 构建分析响应
	response := &dto.StockAnalysisResponse{
		Symbol:      req.Symbol,
		CompanyName: s.extractCompanyName(quote),
		CurrentPrice: s.extractCurrentPrice(quote),
		Currency:    s.extractCurrency(quote),
		LastUpdated: time.Now(),
	}

	// 5. 根据分析类型执行相应分析
	analysisType := req.AnalysisType
	if analysisType == "" {
		analysisType = "all"
	}

	if analysisType == "technical" || analysisType == "all" {
		if history != nil {
			response.TechnicalAnalysis = s.performTechnicalAnalysis(history)
		}
	}

	if analysisType == "fundamental" || analysisType == "all" {
		if companyInfo != nil {
			response.FundamentalAnalysis = s.performFundamentalAnalysis(companyInfo, quote)
		}
	}

	if analysisType == "risk" || analysisType == "all" {
		if history != nil {
			response.RiskAssessment = s.performRiskAssessment(history)
		}
	}

	if analysisType == "all" {
		response.InvestmentAdvice = s.generateInvestmentAdvice(response)
	}

	return response, nil
}

// CompareStocks 对比多只股票
func (s *StockAnalysisService) CompareStocks(ctx context.Context, req *dto.StockCompareRequest) (*dto.StockCompareResponse, error) {
	s.logger.Info("开始对比股票", zap.Strings("symbols", req.Symbols))

	var individual []dto.StockAnalysisResponse
	
	// 分析每只股票
	for _, symbol := range req.Symbols {
		analysisReq := &dto.StockAnalysisRequest{
			Symbol:       symbol,
			Period:       req.Period,
			AnalysisType: "all",
		}
		
		analysis, err := s.AnalyzeStock(ctx, analysisReq)
		if err != nil {
			s.logger.Error("分析股票失败", zap.String("symbol", symbol), zap.Error(err))
			continue
		}
		
		individual = append(individual, *analysis)
	}

	if len(individual) == 0 {
		return nil, fmt.Errorf("没有成功分析任何股票")
	}

	// 执行对比分析
	comparison := s.performStockComparison(individual)
	recommendation := s.generateComparisonRecommendation(individual, comparison)

	return &dto.StockCompareResponse{
		Symbols:        req.Symbols,
		Comparison:     comparison,
		Individual:     individual,
		Recommendation: recommendation,
	}, nil
}

// getStockQuote 获取股票报价
func (s *StockAnalysisService) getStockQuote(ctx context.Context, symbol string) (*dto.MCPExecuteResponse, error) {
	req := &dto.MCPExecuteRequest{
		Name: "yahoo_finance",
		Arguments: map[string]interface{}{
			"action": "quote",
			"symbol": symbol,
		},
	}
	
	return s.mcpClient.ExecuteTool(ctx, req)
}

// getStockHistory 获取股票历史数据
func (s *StockAnalysisService) getStockHistory(ctx context.Context, symbol, period, interval string) (*dto.MCPExecuteResponse, error) {
	req := &dto.MCPExecuteRequest{
		Name: "yahoo_finance",
		Arguments: map[string]interface{}{
			"action":   "history",
			"symbol":   symbol,
			"period":   period,
			"interval": interval,
		},
	}
	
	return s.mcpClient.ExecuteTool(ctx, req)
}

// getStockInfo 获取股票公司信息
func (s *StockAnalysisService) getStockInfo(ctx context.Context, symbol string) (*dto.MCPExecuteResponse, error) {
	req := &dto.MCPExecuteRequest{
		Name: "yahoo_finance",
		Arguments: map[string]interface{}{
			"action": "info",
			"symbol": symbol,
		},
	}
	
	return s.mcpClient.ExecuteTool(ctx, req)
}

// extractCompanyName 从报价数据中提取公司名称
func (s *StockAnalysisService) extractCompanyName(quote *dto.MCPExecuteResponse) string {
	if quote == nil || len(quote.Content) == 0 {
		return ""
	}
	
	content := quote.Content[0].Text
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "公司:") || strings.Contains(line, "Company:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

// extractCurrentPrice 从报价数据中提取当前价格
func (s *StockAnalysisService) extractCurrentPrice(quote *dto.MCPExecuteResponse) float64 {
	if quote == nil || len(quote.Content) == 0 {
		return 0
	}
	
	content := quote.Content[0].Text
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "当前价格:") || strings.Contains(line, "Current Price:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				priceStr := strings.TrimSpace(parts[1])
				priceStr = strings.ReplaceAll(priceStr, "$", "")
				priceStr = strings.ReplaceAll(priceStr, ",", "")
				if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
					return price
				}
			}
		}
	}
	return 0
}

// extractCurrency 从报价数据中提取货币
func (s *StockAnalysisService) extractCurrency(quote *dto.MCPExecuteResponse) string {
	if quote == nil || len(quote.Content) == 0 {
		return "USD"
	}
	
	content := quote.Content[0].Text
	if strings.Contains(content, "货币:") || strings.Contains(content, "Currency:") {
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			if strings.Contains(line, "货币:") || strings.Contains(line, "Currency:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}
	return "USD"
}

// performTechnicalAnalysis 执行技术分析
func (s *StockAnalysisService) performTechnicalAnalysis(history *dto.MCPExecuteResponse) *dto.TechnicalAnalysis {
	if history == nil || len(history.Content) == 0 {
		return nil
	}

	// 解析历史数据
	prices := s.parseHistoricalPrices(history.Content[0].Text)
	if len(prices) < 20 {
		return nil
	}

	// 计算技术指标
	rsi := s.calculateRSI(prices, 14)
	ma5 := s.calculateMA(prices, 5)
	ma10 := s.calculateMA(prices, 10)
	ma20 := s.calculateMA(prices, 20)
	ma50 := s.calculateMA(prices, 50)
	ma200 := s.calculateMA(prices, 200)

	// 确定趋势
	trend := s.determineTrend(prices, ma20)
	
	// 计算支撑位和阻力位
	support, resistance := s.calculateSupportResistance(prices)

	// 生成技术信号
	signals := s.generateTechnicalSignals(prices, rsi, ma5, ma20)

	return &dto.TechnicalAnalysis{
		Trend:      trend,
		Support:    support,
		Resistance: resistance,
		RSI:        rsi,
		MovingAverages: &dto.MovingAverages{
			MA5:   ma5,
			MA10:  ma10,
			MA20:  ma20,
			MA50:  ma50,
			MA200: ma200,
		},
		TechnicalSignals: signals,
	}
}

// performFundamentalAnalysis 执行基本面分析
func (s *StockAnalysisService) performFundamentalAnalysis(info *dto.MCPExecuteResponse, quote *dto.MCPExecuteResponse) *dto.FundamentalAnalysis {
	// 这里应该解析公司信息，提取财务指标
	// 由于Yahoo Finance API的限制，这里提供一个基础实现
	return &dto.FundamentalAnalysis{
		MarketCap:     0, // 需要从API响应中解析
		PE:            0, // 需要从API响应中解析
		PB:            0, // 需要从API响应中解析
		DividendYield: 0, // 需要从API响应中解析
		Valuation:     "需要更多数据",
	}
}

// performRiskAssessment 执行风险评估
func (s *StockAnalysisService) performRiskAssessment(history *dto.MCPExecuteResponse) *dto.RiskAssessment {
	if history == nil || len(history.Content) == 0 {
		return nil
	}

	prices := s.parseHistoricalPrices(history.Content[0].Text)
	if len(prices) < 30 {
		return nil
	}

	// 计算收益率
	returns := s.calculateReturns(prices)
	
	// 计算波动率
	volatility := s.calculateVolatility(returns)
	
	// 计算最大回撤
	maxDrawdown := s.calculateMaxDrawdown(prices)
	
	// 确定风险等级
	riskLevel := s.determineRiskLevel(volatility, maxDrawdown)

	return &dto.RiskAssessment{
		RiskLevel:   riskLevel,
		Volatility:  volatility,
		Beta:        1.0, // 需要市场数据计算
		MaxDrawdown: maxDrawdown,
		RiskFactors: []string{"市场风险", "行业风险", "公司特定风险"},
	}
}

// generateInvestmentAdvice 生成投资建议
func (s *StockAnalysisService) generateInvestmentAdvice(analysis *dto.StockAnalysisResponse) *dto.InvestmentAdvice {
	var score float64 = 0.5 // 基础分数
	var reasons []string
	var risks []string

	// 基于技术分析调整分数
	if analysis.TechnicalAnalysis != nil {
		if analysis.TechnicalAnalysis.Trend == "上升" {
			score += 0.2
			reasons = append(reasons, "技术面显示上升趋势")
		} else if analysis.TechnicalAnalysis.Trend == "下降" {
			score -= 0.2
			risks = append(risks, "技术面显示下降趋势")
		}

		if analysis.TechnicalAnalysis.RSI < 30 {
			score += 0.1
			reasons = append(reasons, "RSI显示超卖状态")
		} else if analysis.TechnicalAnalysis.RSI > 70 {
			score -= 0.1
			risks = append(risks, "RSI显示超买状态")
		}
	}

	// 基于风险评估调整分数
	if analysis.RiskAssessment != nil {
		switch analysis.RiskAssessment.RiskLevel {
		case "低":
			score += 0.1
			reasons = append(reasons, "风险水平较低")
		case "高":
			score -= 0.1
			risks = append(risks, "风险水平较高")
		}
	}

	// 确定推荐操作
	var recommendation string
	if score >= 0.8 {
		recommendation = "强烈买入"
	} else if score >= 0.6 {
		recommendation = "买入"
	} else if score >= 0.4 {
		recommendation = "持有"
	} else if score >= 0.2 {
		recommendation = "卖出"
	} else {
		recommendation = "强烈卖出"
	}

	// 计算目标价格
	targetPrice := analysis.CurrentPrice * (1 + (score-0.5)*0.2)

	return &dto.InvestmentAdvice{
		Recommendation: recommendation,
		TargetPrice:    targetPrice,
		TimeHorizon:    "3-6个月",
		Confidence:     score,
		Reasons:        reasons,
		Risks:          risks,
	}
}

// 辅助函数实现

// parseHistoricalPrices 解析历史价格数据
func (s *StockAnalysisService) parseHistoricalPrices(content string) []float64 {
	var prices []float64
	lines := strings.Split(content, "\n")
	
	for _, line := range lines {
		if strings.Contains(line, "收盘价:") || strings.Contains(line, "Close:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				priceStr := strings.TrimSpace(parts[1])
				priceStr = strings.ReplaceAll(priceStr, "$", "")
				priceStr = strings.ReplaceAll(priceStr, ",", "")
				if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
					prices = append(prices, price)
				}
			}
		}
	}
	
	return prices
}

// calculateRSI 计算RSI指标
func (s *StockAnalysisService) calculateRSI(prices []float64, period int) float64 {
	if len(prices) < period+1 {
		return 50 // 默认值
	}

	gains := 0.0
	losses := 0.0

	// 计算初始平均收益和损失
	for i := 1; i <= period; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return rsi
}

// calculateMA 计算移动平均线
func (s *StockAnalysisService) calculateMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}

	sum := 0.0
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i]
	}

	return sum / float64(period)
}

// determineTrend 确定趋势
func (s *StockAnalysisService) determineTrend(prices []float64, ma20 float64) string {
	if len(prices) == 0 {
		return "横盘"
	}

	currentPrice := prices[len(prices)-1]
	if currentPrice > ma20*1.02 {
		return "上升"
	} else if currentPrice < ma20*0.98 {
		return "下降"
	}
	return "横盘"
}

// calculateSupportResistance 计算支撑位和阻力位
func (s *StockAnalysisService) calculateSupportResistance(prices []float64) (float64, float64) {
	if len(prices) < 20 {
		return 0, 0
	}

	recent := prices[len(prices)-20:]
	sort.Float64s(recent)

	support := recent[len(recent)/4]     // 25%分位数
	resistance := recent[len(recent)*3/4] // 75%分位数

	return support, resistance
}

// generateTechnicalSignals 生成技术信号
func (s *StockAnalysisService) generateTechnicalSignals(prices []float64, rsi, ma5, ma20 float64) []dto.TechnicalSignal {
	var signals []dto.TechnicalSignal

	// RSI信号
	if rsi < 30 {
		signals = append(signals, dto.TechnicalSignal{
			Type:        "RSI",
			Signal:      "买入",
			Strength:    0.8,
			Description: "RSI显示超卖状态",
		})
	} else if rsi > 70 {
		signals = append(signals, dto.TechnicalSignal{
			Type:        "RSI",
			Signal:      "卖出",
			Strength:    0.8,
			Description: "RSI显示超买状态",
		})
	}

	// 移动平均线信号
	if ma5 > ma20 {
		signals = append(signals, dto.TechnicalSignal{
			Type:        "MA",
			Signal:      "买入",
			Strength:    0.6,
			Description: "短期均线上穿长期均线",
		})
	} else if ma5 < ma20 {
		signals = append(signals, dto.TechnicalSignal{
			Type:        "MA",
			Signal:      "卖出",
			Strength:    0.6,
			Description: "短期均线下穿长期均线",
		})
	}

	return signals
}

// calculateReturns 计算收益率
func (s *StockAnalysisService) calculateReturns(prices []float64) []float64 {
	var returns []float64
	for i := 1; i < len(prices); i++ {
		ret := (prices[i] - prices[i-1]) / prices[i-1]
		returns = append(returns, ret)
	}
	return returns
}

// calculateVolatility 计算波动率
func (s *StockAnalysisService) calculateVolatility(returns []float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	// 计算平均收益率
	mean := 0.0
	for _, ret := range returns {
		mean += ret
	}
	mean /= float64(len(returns))

	// 计算方差
	variance := 0.0
	for _, ret := range returns {
		variance += math.Pow(ret-mean, 2)
	}
	variance /= float64(len(returns))

	// 年化波动率
	return math.Sqrt(variance) * math.Sqrt(252) // 252个交易日
}

// calculateMaxDrawdown 计算最大回撤
func (s *StockAnalysisService) calculateMaxDrawdown(prices []float64) float64 {
	if len(prices) == 0 {
		return 0
	}

	maxDrawdown := 0.0
	peak := prices[0]

	for _, price := range prices {
		if price > peak {
			peak = price
		}
		drawdown := (peak - price) / peak
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	return maxDrawdown
}

// determineRiskLevel 确定风险等级
func (s *StockAnalysisService) determineRiskLevel(volatility, maxDrawdown float64) string {
	if volatility > 0.3 || maxDrawdown > 0.2 {
		return "高"
	} else if volatility > 0.15 || maxDrawdown > 0.1 {
		return "中"
	}
	return "低"
}

// performStockComparison 执行股票对比
func (s *StockAnalysisService) performStockComparison(stocks []dto.StockAnalysisResponse) *dto.StockComparison {
	if len(stocks) == 0 {
		return nil
	}

	comparison := &dto.StockComparison{
		Performance: &dto.PerformanceComparison{
			Returns1D: make(map[string]float64),
			Returns1W: make(map[string]float64),
			Returns1M: make(map[string]float64),
		},
		Valuation: &dto.ValuationComparison{
			PE:        make(map[string]float64),
			PB:        make(map[string]float64),
			MarketCap: make(map[string]float64),
		},
		Risk: &dto.RiskComparison{
			Volatility:  make(map[string]float64),
			Beta:        make(map[string]float64),
			MaxDrawdown: make(map[string]float64),
		},
	}

	// 填充对比数据
	for _, stock := range stocks {
		if stock.FundamentalAnalysis != nil {
			comparison.Valuation.PE[stock.Symbol] = stock.FundamentalAnalysis.PE
			comparison.Valuation.PB[stock.Symbol] = stock.FundamentalAnalysis.PB
			comparison.Valuation.MarketCap[stock.Symbol] = stock.FundamentalAnalysis.MarketCap
		}

		if stock.RiskAssessment != nil {
			comparison.Risk.Volatility[stock.Symbol] = stock.RiskAssessment.Volatility
			comparison.Risk.Beta[stock.Symbol] = stock.RiskAssessment.Beta
			comparison.Risk.MaxDrawdown[stock.Symbol] = stock.RiskAssessment.MaxDrawdown
		}
	}

	return comparison
}

// generateComparisonRecommendation 生成对比推荐
func (s *StockAnalysisService) generateComparisonRecommendation(stocks []dto.StockAnalysisResponse, comparison *dto.StockComparison) string {
	if len(stocks) == 0 {
		return "无法生成推荐"
	}

	// 简单的推荐逻辑：选择风险调整后收益最好的股票
	bestStock := stocks[0].Symbol
	bestScore := 0.0

	for _, stock := range stocks {
		score := 0.5 // 基础分数

		if stock.InvestmentAdvice != nil {
			score = stock.InvestmentAdvice.Confidence
		}

		if score > bestScore {
			bestScore = score
			bestStock = stock.Symbol
		}
	}

	return fmt.Sprintf("基于综合分析，推荐关注 %s，其风险调整后的投资价值相对较高", bestStock)
}