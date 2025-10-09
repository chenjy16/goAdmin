package middleware

import (
	"net/http"

	"goMcp/internal/response"

	"github.com/gin-gonic/gin"
)

// Recovery 自定义错误恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			response.InternalServerError(c, "Internal Server Error", err)
		} else {
			response.InternalServerError(c, "Internal Server Error", "Unknown error")
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
