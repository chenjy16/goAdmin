package middleware

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"admin/internal/response"
	"admin/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)





// ValidateJSONFactory 创建JSON验证中间件的工厂函数
func ValidateJSONFactory(objType interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建对象类型的新实例
		obj := createNewInstance(objType)

		if err := c.ShouldBindJSON(obj); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errorMessages := formatValidationErrors(validationErrors)
				response.BadRequest(c, "Validation failed", strings.Join(errorMessages, "; "))
				c.Abort()
				return
			}
			response.BadRequest(c, "Invalid JSON format", err.Error())
			c.Abort()
			return
		}

		// 将验证后的对象存储到上下文中
		c.Set("validated_data", obj)
		c.Next()
	}
}



// formatValidationErrors 格式化验证错误信息
func formatValidationErrors(validationErrors validator.ValidationErrors) []string {
	var errorMessages []string
	for _, e := range validationErrors {
		errorMessages = append(errorMessages, utils.GetValidationErrorMessage(e))
	}
	return errorMessages
}

// getErrorMessage 获取友好的错误信息




// CustomValidationError 自定义验证错误
type CustomValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrorResponse 验证错误响应
type ValidationErrorResponse struct {
	Message string                  `json:"message"`
	Errors  []CustomValidationError `json:"errors"`
}

// HandleValidationError 处理验证错误的通用函数
func HandleValidationError(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errors []CustomValidationError

		for _, e := range validationErrors {
			errors = append(errors, CustomValidationError{
				Field:   e.Field(),
				Message: utils.GetValidationErrorMessage(e),
				Value:   fmt.Sprintf("%v", e.Value()),
			})
		}

		response := ValidationErrorResponse{
			Message: "Validation failed",
			Errors:  errors,
		}

		c.JSON(http.StatusBadRequest, response)
		return
	}

	response.BadRequest(c, "Invalid request", err.Error())
}

// createNewInstance 创建对象类型的新实例
func createNewInstance(objType interface{}) interface{} {
	// 使用反射创建新实例
	t := reflect.TypeOf(objType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t).Interface()
}
