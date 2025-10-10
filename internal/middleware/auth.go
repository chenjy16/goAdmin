package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"go-springAi/internal/errors"
	"go-springAi/internal/response"
	"go-springAi/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(jwtManager *utils.JWTManager, zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			zapLogger.Warn("Unauthorized access attempt",
				zap.String("module", "auth"),
				zap.String("component", "middleware"),
				zap.String("operation", "auth"),
				zap.String("reason", "missing_authorization_header"))
			
			response.Error(c, http.StatusUnauthorized, "Authorization header required", "")
			c.Abort()
			return
		}

		// 检查Bearer token格式
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			zapLogger.Warn("Unauthorized access attempt",
				zap.String("module", "auth"),
				zap.String("component", "middleware"),
				zap.String("operation", "auth"),
				zap.String("reason", "invalid_token_format"))
			
			response.Error(c, http.StatusUnauthorized, "Invalid token format", "")
			c.Abort()
			return
		}

		token := tokenParts[1]

		// 验证JWT token
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			zapLogger.Warn("Token validation failed",
				zap.String("module", "auth"),
				zap.String("component", "middleware"),
				zap.String("operation", "auth"),
				zap.Error(err))
			
			// 根据错误类型返回不同的响应
			if strings.Contains(err.Error(), "expired") {
				response.Error(c, http.StatusUnauthorized, "Token expired", "")
			} else {
				response.Error(c, http.StatusUnauthorized, "Invalid token", "")
			}
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", strconv.FormatInt(claims.UserID, 10))
		c.Set("username", claims.Username)
		c.Set("auth_type", "jwt")

		zapLogger.Info("Token validated successfully",
			zap.String("module", "auth"),
			zap.String("component", "middleware"),
			zap.String("operation", "auth"),
			zap.String("user_id", strconv.FormatInt(claims.UserID, 10)),
			zap.String("username", claims.Username))

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件
// 如果提供了token则验证，如果没有提供则继续执行但不设置用户信息
func OptionalAuthMiddleware(jwtManager *utils.JWTManager, zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有提供认证头，继续执行
			c.Next()
			return
		}

		// 检查Bearer token格式
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			// token格式错误，继续执行但不设置用户信息
			c.Next()
			return
		}

		token := tokenParts[1]

		// 验证JWT token
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			// token无效，继续执行但不设置用户信息
			zapLogger.Warn("Token validation failed in optional auth",
				zap.String("module", "auth"),
				zap.String("component", "middleware"),
				zap.String("operation", "optional_auth"),
				zap.Error(err))
			c.Next()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", strconv.FormatInt(claims.UserID, 10))
		c.Set("username", claims.Username)
		c.Set("auth_type", "jwt")

		zapLogger.Info("Token validated successfully in optional auth",
			zap.String("module", "auth"),
			zap.String("component", "middleware"),
			zap.String("operation", "optional_auth"),
			zap.String("user_id", strconv.FormatInt(claims.UserID, 10)),
			zap.String("username", claims.Username))

		c.Next()
	}
}

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(c *gin.Context) (int64, error) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return 0, errors.NewUnauthorizedError("User not authenticated")
	}

	userIDString, ok := userIDStr.(string)
	if !ok {
		return 0, errors.NewInternalError("Invalid user ID format in context")
	}

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return 0, errors.NewInternalError("Failed to parse user ID")
	}

	return userID, nil
}

// GetUsernameFromContext 从上下文中获取用户名
func GetUsernameFromContext(c *gin.Context) (string, error) {
	username, exists := c.Get("username")
	if !exists {
		return "", errors.NewUnauthorizedError("User not authenticated")
	}

	usernameStr, ok := username.(string)
	if !ok {
		return "", errors.NewInternalError("Invalid username format in context")
	}

	return usernameStr, nil
}