package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// CustomValidator 自定义验证器
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator 创建自定义验证器实例
func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// 注册自定义验证规则
	v.RegisterValidation("strong_password", validateStrongPassword)
	v.RegisterValidation("username_format", validateUsernameFormat)
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("chinese_name", validateChineseName)

	// 注册字段名标签函数
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{validator: v}
}

// ValidateStruct 验证结构体
func (cv *CustomValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		if err := cv.validator.Struct(obj); err != nil {
			return cv.formatValidationError(err)
		}
	}
	return nil
}

// Engine 返回验证器引擎
func (cv *CustomValidator) Engine() interface{} {
	return cv.validator
}

// formatValidationError 格式化验证错误信息
func (cv *CustomValidator) formatValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, e := range validationErrors {
			errorMessages = append(errorMessages, cv.getErrorMessage(e))
		}
		return fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))
	}
	return err
}

// getErrorMessage 获取友好的错误信息
// GetValidationErrorMessage 获取验证错误的友好消息（公共函数）
func GetValidationErrorMessage(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, e.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, e.Param())
	case "numeric":
		return fmt.Sprintf("%s must be a number", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "strong_password":
		return fmt.Sprintf("%s must contain at least 8 characters with uppercase, lowercase, number and special character", field)
	case "username_format":
		return fmt.Sprintf("%s must be 3-20 characters long and contain only letters, numbers, and underscores", field)
	case "phone":
		return fmt.Sprintf("%s must be a valid phone number", field)
	case "chinese_name":
		return fmt.Sprintf("%s must be a valid Chinese name", field)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, e.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func (cv *CustomValidator) getErrorMessage(e validator.FieldError) string {
	return GetValidationErrorMessage(e)
}

// kindOfData 获取数据类型
func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

// 自定义验证函数

// validateStrongPassword 验证强密码
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateUsernameFormat 验证用户名格式
func validateUsernameFormat(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	if len(username) < 3 || len(username) > 50 {
		return false
	}

	for _, char := range username {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_') {
			return false
		}
	}

	return true
}

// validatePhone 验证手机号
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// 使用正则表达式验证中国手机号
	if len(phone) != 11 {
		return false
	}

	// 简单的中国手机号验证
	if phone[0] != '1' {
		return false
	}

	for _, char := range phone {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// validateChineseName 验证中文姓名
func validateChineseName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	if len(name) < 2 || len(name) > 20 {
		return false
	}

	// 简单的中文字符验证
	for _, char := range name {
		if char < 0x4e00 || char > 0x9fff {
			return false
		}
	}

	return true
}
