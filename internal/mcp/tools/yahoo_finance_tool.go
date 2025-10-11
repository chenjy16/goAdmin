package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-springAi/internal/dto"
	"go-springAi/internal/mcp"
)

// YahooFinanceTool Yahoo Finance 股票数据工具
type YahooFinanceTool struct {
	*mcp.BaseTool
	httpClient *http.Client
}

// NewYahooFinanceTool 创建 Yahoo Finance 工具
func NewYahooFinanceTool() *YahooFinanceTool {
	return &YahooFinanceTool{
		BaseTool: &mcp.BaseTool{
			Name:        "雅虎财经",
			Description: "获取股票数据",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"action": map[string]interface{}{
						"type":        "string",
						"description": "Action to perform: 'quote', 'history', or 'info'",
						"enum":        []string{"quote", "history", "info"},
					},
					"symbol": map[string]interface{}{
						"type":        "string",
						"description": "Stock symbol (e.g., AAPL, TSLA, MSFT)",
					},
					"period": map[string]interface{}{
						"type":        "string",
						"description": "Time period for historical data: '1d', '5d', '1mo', '3mo', '6mo', '1y', '2y', '5y', '10y', 'ytd', 'max'",
						"enum":        []string{"1d", "5d", "1mo", "3mo", "6mo", "1y", "2y", "5y", "10y", "ytd", "max"},
						"default":     "1mo",
					},
					"interval": map[string]interface{}{
						"type":        "string",
						"description": "Data interval: '1m', '2m', '5m', '15m', '30m', '60m', '90m', '1h', '1d', '5d', '1wk', '1mo', '3mo'",
						"enum":        []string{"1m", "2m", "5m", "15m", "30m", "60m", "90m", "1h", "1d", "5d", "1wk", "1mo", "3mo"},
						"default":     "1d",
					},
				},
				"required": []string{"action", "symbol"},
			},
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Execute 执行 Yahoo Finance 工具
func (yf *YahooFinanceTool) Execute(ctx context.Context, args map[string]interface{}) (*dto.MCPExecuteResponse, error) {
	// 验证参数
	if err := yf.Validate(args); err != nil {
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

	action := args["action"].(string)
	symbol := strings.ToUpper(args["symbol"].(string))

	switch action {
	case "quote":
		return yf.getQuote(ctx, symbol)
	case "history":
		period := "1mo"
		interval := "1d"
		if p, ok := args["period"].(string); ok {
			period = p
		}
		if i, ok := args["interval"].(string); ok {
			interval = i
		}
		return yf.getHistory(ctx, symbol, period, interval)
	case "info":
		return yf.getInfo(ctx, symbol)
	default:
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("不支持的操作: %s", action),
				},
			},
			IsError: true,
		}, nil
	}
}

// Validate 验证参数
func (yf *YahooFinanceTool) Validate(args map[string]interface{}) error {
	// 检查必需参数
	action, ok := args["action"].(string)
	if !ok {
		return fmt.Errorf("action 参数是必需的且必须是字符串")
	}

	symbol, ok := args["symbol"].(string)
	if !ok {
		return fmt.Errorf("symbol 参数是必需的且必须是字符串")
	}

	if symbol == "" {
		return fmt.Errorf("symbol 不能为空")
	}

	// 验证 action 值
	validActions := []string{"quote", "history", "info"}
	actionValid := false
	for _, validAction := range validActions {
		if action == validAction {
			actionValid = true
			break
		}
	}
	if !actionValid {
		return fmt.Errorf("action 必须是以下值之一: %v", validActions)
	}

	return nil
}

// getQuote 获取股票实时报价
func (yf *YahooFinanceTool) getQuote(ctx context.Context, symbol string) (*dto.MCPExecuteResponse, error) {
	// 使用 Yahoo Finance API v8
	apiURL := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s", symbol)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("创建请求失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := yf.httpClient.Do(req)
	if err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("请求失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("读取响应失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// 解析 Yahoo Finance 响应
	var yahooResp YahooFinanceResponse
	if err := json.Unmarshal(body, &yahooResp); err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("解析响应失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	if yahooResp.Chart.Error != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Yahoo Finance API 错误: %s", yahooResp.Chart.Error.Description),
				},
			},
			IsError: true,
		}, nil
	}

	if len(yahooResp.Chart.Result) == 0 {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("未找到股票 %s 的数据", symbol),
				},
			},
			IsError: true,
		}, nil
	}

	result := yahooResp.Chart.Result[0]
	meta := result.Meta

	// 格式化股票报价信息
	quoteText := fmt.Sprintf(`📈 %s (%s) 股票报价

💰 当前价格: $%.2f
📊 前收盘价: $%.2f
📈 今日开盘: $%.2f
🔺 今日最高: $%.2f
🔻 今日最低: $%.2f
📊 成交量: %s
🏢 市场: %s
💱 货币: %s
⏰ 更新时间: %s`,
		meta.Symbol,
		meta.Symbol,
		meta.RegularMarketPrice,
		meta.PreviousClose,
		meta.RegularMarketDayHigh,
		meta.RegularMarketDayHigh,
		meta.RegularMarketDayLow,
		formatVolume(meta.RegularMarketVolume),
		meta.ExchangeName,
		meta.Currency,
		time.Unix(meta.RegularMarketTime, 0).Format("2006-01-02 15:04:05"))

	// 计算涨跌幅
	if meta.PreviousClose > 0 {
		change := meta.RegularMarketPrice - meta.PreviousClose
		changePercent := (change / meta.PreviousClose) * 100
		changeEmoji := "📈"
		if change < 0 {
			changeEmoji = "📉"
		}
		quoteText += fmt.Sprintf("\n%s 涨跌: $%.2f (%.2f%%)", changeEmoji, change, changePercent)
	}

	return &dto.MCPExecuteResponse{
		Content: []dto.MCPContent{
			{
				Type: "text",
				Text: quoteText,
			},
		},
		IsError: false,
	}, nil
}

// getHistory 获取股票历史数据
func (yf *YahooFinanceTool) getHistory(ctx context.Context, symbol, period, interval string) (*dto.MCPExecuteResponse, error) {
	// 构建 URL 参数
	params := url.Values{}
	params.Set("period1", "0")
	params.Set("period2", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("interval", interval)
	params.Set("includePrePost", "true")
	params.Set("events", "div,splits")

	// 根据 period 计算开始时间
	now := time.Now()
	var startTime time.Time
	switch period {
	case "1d":
		startTime = now.AddDate(0, 0, -1)
	case "5d":
		startTime = now.AddDate(0, 0, -5)
	case "1mo":
		startTime = now.AddDate(0, -1, 0)
	case "3mo":
		startTime = now.AddDate(0, -3, 0)
	case "6mo":
		startTime = now.AddDate(0, -6, 0)
	case "1y":
		startTime = now.AddDate(-1, 0, 0)
	case "2y":
		startTime = now.AddDate(-2, 0, 0)
	case "5y":
		startTime = now.AddDate(-5, 0, 0)
	case "10y":
		startTime = now.AddDate(-10, 0, 0)
	case "ytd":
		startTime = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	case "max":
		startTime = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	default:
		startTime = now.AddDate(0, -1, 0) // 默认1个月
	}

	params.Set("period1", strconv.FormatInt(startTime.Unix(), 10))

	apiURL := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?%s", symbol, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("创建请求失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := yf.httpClient.Do(req)
	if err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("请求失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("读取响应失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	var yahooResp YahooFinanceResponse
	if err := json.Unmarshal(body, &yahooResp); err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("解析响应失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	if yahooResp.Chart.Error != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Yahoo Finance API 错误: %s", yahooResp.Chart.Error.Description),
				},
			},
			IsError: true,
		}, nil
	}

	if len(yahooResp.Chart.Result) == 0 {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("未找到股票 %s 的历史数据", symbol),
				},
			},
			IsError: true,
		}, nil
	}

	result := yahooResp.Chart.Result[0]

	// 格式化历史数据
	historyText := fmt.Sprintf("📊 %s 历史数据 (%s, %s)\n\n", symbol, period, interval)

	if len(result.Timestamp) > 0 && result.Indicators.Quote != nil && len(result.Indicators.Quote) > 0 {
		quote := result.Indicators.Quote[0]

		// 显示最近几个数据点
		maxPoints := 10
		if len(result.Timestamp) < maxPoints {
			maxPoints = len(result.Timestamp)
		}

		for i := len(result.Timestamp) - maxPoints; i < len(result.Timestamp); i++ {
			timestamp := time.Unix(result.Timestamp[i], 0)

			if i < len(quote.Open) && i < len(quote.High) && i < len(quote.Low) && i < len(quote.Close) && i < len(quote.Volume) {
				historyText += fmt.Sprintf("📅 %s\n", timestamp.Format("2006-01-02 15:04"))
				historyText += fmt.Sprintf("   开盘: $%.2f | 最高: $%.2f | 最低: $%.2f | 收盘: $%.2f\n",
					quote.Open[i], quote.High[i], quote.Low[i], quote.Close[i])
				historyText += fmt.Sprintf("   成交量: %s\n\n", formatVolume(int64(quote.Volume[i])))
			}
		}
	}

	return &dto.MCPExecuteResponse{
		Content: []dto.MCPContent{
			{
				Type: "text",
				Text: historyText,
			},
		},
		IsError: false,
	}, nil
}

// getInfo 获取股票基本信息
func (yf *YahooFinanceTool) getInfo(ctx context.Context, symbol string) (*dto.MCPExecuteResponse, error) {
	// 使用 Yahoo Finance quoteSummary API
	modules := []string{"summaryProfile", "summaryDetail", "financialData", "defaultKeyStatistics"}
	apiURL := fmt.Sprintf("https://query1.finance.yahoo.com/v10/finance/quoteSummary/%s?modules=%s",
		symbol, strings.Join(modules, ","))

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("创建请求失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := yf.httpClient.Do(req)
	if err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("请求失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("读取响应失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	var summaryResp YahooSummaryResponse
	if err := json.Unmarshal(body, &summaryResp); err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("解析响应失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	if summaryResp.QuoteSummary.Error != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Yahoo Finance API 错误: %s", summaryResp.QuoteSummary.Error.Description),
				},
			},
			IsError: true,
		}, nil
	}

	if len(summaryResp.QuoteSummary.Result) == 0 {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("未找到股票 %s 的公司信息", symbol),
				},
			},
			IsError: true,
		}, nil
	}

	result := summaryResp.QuoteSummary.Result[0]

	// 格式化公司信息
	infoText := fmt.Sprintf("🏢 %s 公司信息\n\n", symbol)

	if result.SummaryProfile != nil {
		profile := result.SummaryProfile
		infoText += fmt.Sprintf("📝 公司名称: %s\n", profile.LongName)
		infoText += fmt.Sprintf("🏭 行业: %s\n", profile.Industry)
		infoText += fmt.Sprintf("🏢 板块: %s\n", profile.Sector)
		infoText += fmt.Sprintf("🌍 国家: %s\n", profile.Country)
		infoText += fmt.Sprintf("🌐 网站: %s\n", profile.Website)
		infoText += fmt.Sprintf("👥 员工数: %s\n", formatNumber(profile.FullTimeEmployees))
		if profile.LongBusinessSummary != "" {
			summary := profile.LongBusinessSummary
			if len(summary) > 200 {
				summary = summary[:200] + "..."
			}
			infoText += fmt.Sprintf("📄 业务简介: %s\n", summary)
		}
		infoText += "\n"
	}

	if result.SummaryDetail != nil {
		detail := result.SummaryDetail
		infoText += "📊 关键指标:\n"
		if detail.MarketCap != nil {
			infoText += fmt.Sprintf("💰 市值: $%s\n", formatLargeNumber(detail.MarketCap.Raw))
		}
		if detail.PeRatio != nil {
			infoText += fmt.Sprintf("📈 市盈率: %.2f\n", detail.PeRatio.Raw)
		}
		if detail.DividendYield != nil {
			infoText += fmt.Sprintf("💵 股息收益率: %.2f%%\n", detail.DividendYield.Raw*100)
		}
		if detail.Beta != nil {
			infoText += fmt.Sprintf("📊 Beta: %.2f\n", detail.Beta.Raw)
		}
	}

	return &dto.MCPExecuteResponse{
		Content: []dto.MCPContent{
			{
				Type: "text",
				Text: infoText,
			},
		},
		IsError: false,
	}, nil
}

// formatVolume 格式化成交量
func formatVolume(volume int64) string {
	if volume >= 1000000000 {
		return fmt.Sprintf("%.1fB", float64(volume)/1000000000)
	} else if volume >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(volume)/1000000)
	} else if volume >= 1000 {
		return fmt.Sprintf("%.1fK", float64(volume)/1000)
	}
	return fmt.Sprintf("%d", volume)
}

// formatNumber 格式化数字
func formatNumber(num int64) string {
	if num >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(num)/1000000)
	} else if num >= 1000 {
		return fmt.Sprintf("%.1fK", float64(num)/1000)
	}
	return fmt.Sprintf("%d", num)
}

// formatLargeNumber 格式化大数字
func formatLargeNumber(num float64) string {
	if num >= 1000000000000 {
		return fmt.Sprintf("%.2fT", num/1000000000000)
	} else if num >= 1000000000 {
		return fmt.Sprintf("%.2fB", num/1000000000)
	} else if num >= 1000000 {
		return fmt.Sprintf("%.2fM", num/1000000)
	}
	return fmt.Sprintf("%.0f", num)
}

// Yahoo Finance API 响应结构体
type YahooFinanceResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency             string  `json:"currency"`
				Symbol               string  `json:"symbol"`
				ExchangeName         string  `json:"exchangeName"`
				RegularMarketPrice   float64 `json:"regularMarketPrice"`
				PreviousClose        float64 `json:"previousClose"`
				RegularMarketDayHigh float64 `json:"regularMarketDayHigh"`
				RegularMarketDayLow  float64 `json:"regularMarketDayLow"`
				RegularMarketVolume  int64   `json:"regularMarketVolume"`
				RegularMarketTime    int64   `json:"regularMarketTime"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []float64 `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"chart"`
}

// Yahoo Finance Summary API 响应结构体
type YahooSummaryResponse struct {
	QuoteSummary struct {
		Result []struct {
			SummaryProfile *struct {
				LongName            string `json:"longName"`
				Industry            string `json:"industry"`
				Sector              string `json:"sector"`
				Country             string `json:"country"`
				Website             string `json:"website"`
				FullTimeEmployees   int64  `json:"fullTimeEmployees"`
				LongBusinessSummary string `json:"longBusinessSummary"`
			} `json:"summaryProfile"`
			SummaryDetail *struct {
				MarketCap *struct {
					Raw float64 `json:"raw"`
				} `json:"marketCap"`
				PeRatio *struct {
					Raw float64 `json:"raw"`
				} `json:"trailingPE"`
				DividendYield *struct {
					Raw float64 `json:"raw"`
				} `json:"dividendYield"`
				Beta *struct {
					Raw float64 `json:"raw"`
				} `json:"beta"`
			} `json:"summaryDetail"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"quoteSummary"`
}
