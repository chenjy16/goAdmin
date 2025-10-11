package handler

import (
	"context"
	"net/http"

	"go-springAi/internal/dto"
	"go-springAi/internal/response"
	"go-springAi/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// StockHandler 股票处理器
type StockHandler struct {
	stockAnalysisService *service.StockAnalysisService
	logger               *zap.Logger
}

// NewStockHandler 创建新的股票处理器
func NewStockHandler(stockAnalysisService *service.StockAnalysisService, logger *zap.Logger) *StockHandler {
	return &StockHandler{
		stockAnalysisService: stockAnalysisService,
		logger:               logger,
	}
}

// AnalyzeStock 分析股票
func (h *StockHandler) AnalyzeStock(c *gin.Context) {
	var req dto.StockAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("绑定股票分析请求失败", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "请求参数无效", err.Error())
		return
	}

	// 验证必需字段
	if req.Symbol == "" {
		response.Error(c, http.StatusBadRequest, "股票代码不能为空", "")
		return
	}

	// 调用股票分析服务
	result, err := h.stockAnalysisService.AnalyzeStock(context.Background(), &req)
	if err != nil {
		h.logger.Error("股票分析失败", zap.Error(err), zap.String("symbol", req.Symbol))
		response.Error(c, http.StatusInternalServerError, "股票分析失败", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "股票分析成功", result)
}

// CompareStocks 对比股票
func (h *StockHandler) CompareStocks(c *gin.Context) {
	var req dto.StockCompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("绑定股票对比请求失败", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "请求参数无效", err.Error())
		return
	}

	// 验证必需字段
	if len(req.Symbols) < 2 {
		response.Error(c, http.StatusBadRequest, "至少需要两个股票代码进行对比", "")
		return
	}

	// 调用股票对比服务
	result, err := h.stockAnalysisService.CompareStocks(context.Background(), &req)
	if err != nil {
		h.logger.Error("股票对比失败", zap.Error(err), zap.Strings("symbols", req.Symbols))
		response.Error(c, http.StatusInternalServerError, "股票对比失败", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "股票对比成功", result)
}

// GetStockQuote 获取股票报价
func (h *StockHandler) GetStockQuote(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		response.Error(c, http.StatusBadRequest, "股票代码不能为空", "")
		return
	}

	// 创建简单的分析请求
	req := &dto.StockAnalysisRequest{
		Symbol:       symbol,
		Period:       "1d",
		AnalysisType: "basic",
	}

	// 调用股票分析服务获取基本信息
	result, err := h.stockAnalysisService.AnalyzeStock(context.Background(), req)
	if err != nil {
		h.logger.Error("获取股票报价失败", zap.Error(err), zap.String("symbol", symbol))
		response.Error(c, http.StatusInternalServerError, "获取股票报价失败", err.Error())
		return
	}

	// 返回简化的报价信息
	quote := map[string]interface{}{
		"symbol":        result.Symbol,
		"current_price": result.CurrentPrice,
		"currency":      result.Currency,
		"company_name":  result.CompanyName,
	}

	response.Success(c, http.StatusOK, "获取股票报价成功", quote)
}

// GetStockHistory 获取股票历史数据
func (h *StockHandler) GetStockHistory(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		response.Error(c, http.StatusBadRequest, "股票代码不能为空", "")
		return
	}

	// 获取查询参数
	period := c.DefaultQuery("period", "1y")
	analysisType := c.DefaultQuery("analysis_type", "technical")

	// 验证参数
	validPeriods := map[string]bool{
		"1d": true, "5d": true, "1mo": true, "3mo": true,
		"6mo": true, "1y": true, "2y": true, "5y": true, "10y": true,
	}
	if !validPeriods[period] {
		response.Error(c, http.StatusBadRequest, "无效的时间周期", "")
		return
	}

	// 创建分析请求
	req := &dto.StockAnalysisRequest{
		Symbol:       symbol,
		Period:       period,
		AnalysisType: analysisType,
	}

	// 调用股票分析服务
	result, err := h.stockAnalysisService.AnalyzeStock(context.Background(), req)
	if err != nil {
		h.logger.Error("获取股票历史数据失败", zap.Error(err), zap.String("symbol", symbol))
		response.Error(c, http.StatusInternalServerError, "获取股票历史数据失败", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "获取股票历史数据成功", result)
}

// GetMarketSummary 获取市场概览
func (h *StockHandler) GetMarketSummary(c *gin.Context) {
	// 主要市场指数
	majorIndices := []string{"^GSPC", "^DJI", "^IXIC", "^RUT"}
	
	var summaries []map[string]interface{}
	
	for _, symbol := range majorIndices {
		req := &dto.StockAnalysisRequest{
			Symbol:       symbol,
			Period:       "1d",
			AnalysisType: "basic",
		}
		
		result, err := h.stockAnalysisService.AnalyzeStock(context.Background(), req)
		if err != nil {
			h.logger.Warn("获取市场指数失败", zap.Error(err), zap.String("symbol", symbol))
			continue
		}
		
		summary := map[string]interface{}{
			"symbol":        result.Symbol,
			"current_price": result.CurrentPrice,
			"currency":      result.Currency,
		}
		summaries = append(summaries, summary)
	}
	
	response.Success(c, http.StatusOK, "获取市场概览成功", map[string]interface{}{
		"indices": summaries,
	})
}