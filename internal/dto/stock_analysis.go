package dto

import "time"

// StockAnalysisRequest 股票分析请求
type StockAnalysisRequest struct {
	Symbol     string `json:"symbol" binding:"required"`     // 股票代码
	Period     string `json:"period,omitempty"`              // 分析周期 (1d, 5d, 1mo, 3mo, 6mo, 1y, 2y, 5y, 10y, ytd, max)
	AnalysisType string `json:"analysis_type,omitempty"`     // 分析类型 (technical, fundamental, risk, all)
}

// StockCompareRequest 股票对比请求
type StockCompareRequest struct {
	Symbols []string `json:"symbols" binding:"required,min=2,max=5"` // 要对比的股票代码列表
	Period  string   `json:"period,omitempty"`                       // 对比周期
}

// StockAnalysisResponse 股票分析响应
type StockAnalysisResponse struct {
	Symbol           string                `json:"symbol"`
	CompanyName      string                `json:"company_name"`
	CurrentPrice     float64               `json:"current_price"`
	Currency         string                `json:"currency"`
	LastUpdated      time.Time             `json:"last_updated"`
	TechnicalAnalysis *TechnicalAnalysis   `json:"technical_analysis,omitempty"`
	FundamentalAnalysis *FundamentalAnalysis `json:"fundamental_analysis,omitempty"`
	RiskAssessment   *RiskAssessment       `json:"risk_assessment,omitempty"`
	InvestmentAdvice *InvestmentAdvice     `json:"investment_advice,omitempty"`
}

// TechnicalAnalysis 技术分析
type TechnicalAnalysis struct {
	Trend            string             `json:"trend"`              // 趋势 (上升/下降/横盘)
	Support          float64            `json:"support"`            // 支撑位
	Resistance       float64            `json:"resistance"`         // 阻力位
	RSI              float64            `json:"rsi"`                // 相对强弱指数
	MACD             *MACDIndicator     `json:"macd,omitempty"`     // MACD指标
	MovingAverages   *MovingAverages    `json:"moving_averages,omitempty"` // 移动平均线
	TechnicalSignals []TechnicalSignal  `json:"technical_signals,omitempty"` // 技术信号
}

// MACDIndicator MACD指标
type MACDIndicator struct {
	MACD      float64 `json:"macd"`
	Signal    float64 `json:"signal"`
	Histogram float64 `json:"histogram"`
	Signal_   string  `json:"signal_type"` // 买入/卖出/持有
}

// MovingAverages 移动平均线
type MovingAverages struct {
	MA5   float64 `json:"ma5"`   // 5日移动平均
	MA10  float64 `json:"ma10"`  // 10日移动平均
	MA20  float64 `json:"ma20"`  // 20日移动平均
	MA50  float64 `json:"ma50"`  // 50日移动平均
	MA200 float64 `json:"ma200"` // 200日移动平均
}

// TechnicalSignal 技术信号
type TechnicalSignal struct {
	Type        string  `json:"type"`        // 信号类型
	Signal      string  `json:"signal"`      // 买入/卖出/持有
	Strength    float64 `json:"strength"`    // 信号强度 (0-1)
	Description string  `json:"description"` // 信号描述
}

// FundamentalAnalysis 基本面分析
type FundamentalAnalysis struct {
	MarketCap       float64 `json:"market_cap"`        // 市值
	PE              float64 `json:"pe"`                // 市盈率
	PB              float64 `json:"pb"`                // 市净率
	DividendYield   float64 `json:"dividend_yield"`    // 股息率
	ROE             float64 `json:"roe"`               // 净资产收益率
	DebtToEquity    float64 `json:"debt_to_equity"`    // 负债权益比
	RevenueGrowth   float64 `json:"revenue_growth"`    // 营收增长率
	EarningsGrowth  float64 `json:"earnings_growth"`   // 盈利增长率
	Valuation       string  `json:"valuation"`         // 估值水平 (低估/合理/高估)
}

// RiskAssessment 风险评估
type RiskAssessment struct {
	RiskLevel    string  `json:"risk_level"`    // 风险等级 (低/中/高)
	Volatility   float64 `json:"volatility"`    // 波动率
	Beta         float64 `json:"beta"`          // 贝塔系数
	MaxDrawdown  float64 `json:"max_drawdown"`  // 最大回撤
	VaR          float64 `json:"var"`           // 风险价值
	RiskFactors  []string `json:"risk_factors"` // 风险因素
}

// InvestmentAdvice 投资建议
type InvestmentAdvice struct {
	Recommendation string   `json:"recommendation"` // 推荐操作 (强烈买入/买入/持有/卖出/强烈卖出)
	TargetPrice    float64  `json:"target_price"`   // 目标价格
	TimeHorizon    string   `json:"time_horizon"`   // 投资时间范围
	Confidence     float64  `json:"confidence"`     // 建议置信度 (0-1)
	Reasons        []string `json:"reasons"`        // 建议理由
	Risks          []string `json:"risks"`          // 潜在风险
}

// StockCompareResponse 股票对比响应
type StockCompareResponse struct {
	Symbols     []string                `json:"symbols"`
	Comparison  *StockComparison        `json:"comparison"`
	Individual  []StockAnalysisResponse `json:"individual"`
	Recommendation string               `json:"recommendation"` // 对比后的推荐
}

// StockComparison 股票对比
type StockComparison struct {
	Performance *PerformanceComparison `json:"performance,omitempty"`
	Valuation   *ValuationComparison   `json:"valuation,omitempty"`
	Risk        *RiskComparison        `json:"risk,omitempty"`
}

// PerformanceComparison 表现对比
type PerformanceComparison struct {
	Returns1D  map[string]float64 `json:"returns_1d"`  // 1日收益率
	Returns1W  map[string]float64 `json:"returns_1w"`  // 1周收益率
	Returns1M  map[string]float64 `json:"returns_1m"`  // 1月收益率
	Returns3M  map[string]float64 `json:"returns_3m"`  // 3月收益率
	Returns1Y  map[string]float64 `json:"returns_1y"`  // 1年收益率
	BestPerformer string           `json:"best_performer"` // 最佳表现者
}

// ValuationComparison 估值对比
type ValuationComparison struct {
	PE           map[string]float64 `json:"pe"`
	PB           map[string]float64 `json:"pb"`
	MarketCap    map[string]float64 `json:"market_cap"`
	MostUndervalued string          `json:"most_undervalued"` // 最被低估的
}

// RiskComparison 风险对比
type RiskComparison struct {
	Volatility   map[string]float64 `json:"volatility"`
	Beta         map[string]float64 `json:"beta"`
	MaxDrawdown  map[string]float64 `json:"max_drawdown"`
	LowestRisk   string             `json:"lowest_risk"` // 风险最低的
}