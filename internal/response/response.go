package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string, err string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Error:   err,
	})
}

// BadRequest 400错误
func BadRequest(c *gin.Context, message string, err string) {
	Error(c, http.StatusBadRequest, message, err)
}



// InternalServerError 500错误
func InternalServerError(c *gin.Context, message string, err string) {
	Error(c, http.StatusInternalServerError, message, err)
}
