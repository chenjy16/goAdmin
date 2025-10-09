package controllers

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"goMcp/internal/errors"
	"goMcp/internal/response"
	"goMcp/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BaseController 基础控制器，包含公共方法
type BaseController struct{}

// NewBaseController 创建基础控制器实例
func NewBaseController() *BaseController {
	return &BaseController{}
}

// HandleValidationError 处理验证错误，提供更友好的错误信息
func (bc *BaseController) HandleValidationError(c *gin.Context, err error) {
	// 使用验证中间件的错误格式化功能
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, e := range validationErrors {
			errorMessages = append(errorMessages, utils.GetValidationErrorMessage(e))
		}
		response.BadRequest(c, "Validation failed", strings.Join(errorMessages, "; "))
		return
	}

	// 其他类型的绑定错误
	response.BadRequest(c, "Invalid request data", err.Error())
}



// BindAndValidate 统一的数据绑定和验证函数
func (bc *BaseController) BindAndValidate(c *gin.Context, req interface{}) error {
	// 首先尝试从验证中间件获取已验证的数据
	if validatedData, exists := c.Get("validated_data"); exists {
		// 使用反射将验证中间件的数据复制到目标结构体
		return bc.CopyValidatedData(validatedData, req)
	}

	// 如果没有验证中间件，则使用传统的绑定和验证方式
	if err := c.ShouldBindJSON(req); err != nil {
		bc.HandleValidationError(c, err)
		return err
	}

	return nil
}

// CopyValidatedData 将验证中间件的数据复制到目标结构体
func (bc *BaseController) CopyValidatedData(source interface{}, target interface{}) error {
	// 使用反射进行通用的数据复制
	sourceValue := reflect.ValueOf(source)
	targetValue := reflect.ValueOf(target)

	// 检查 target 是否为指针
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	// 获取指针指向的元素
	targetElem := targetValue.Elem()

	// 如果 source 是指针，获取其指向的元素
	if sourceValue.Kind() == reflect.Ptr {
		sourceValue = sourceValue.Elem()
	}

	// 检查类型是否兼容
	if !sourceValue.Type().AssignableTo(targetElem.Type()) {
		return fmt.Errorf("type mismatch: cannot copy %T to %T", source, target)
	}

	// 执行复制
	targetElem.Set(sourceValue)
	return nil
}

// ParseIDParam 解析路径参数中的ID
func (bc *BaseController) ParseIDParam(c *gin.Context, paramName string) (int64, error) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.NewValidationError(fmt.Sprintf("Invalid %s", paramName)).WithDetails(err.Error())
	}
	return id, nil
}

// ParsePaginationParams 解析分页参数
func (bc *BaseController) ParsePaginationParams(c *gin.Context) (page, limit, offset int64, err error) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		err = errors.NewValidationError("Invalid page number").WithDetails("Page must be a positive integer")
		return
	}

	limit, err = strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit < 1 || limit > 100 {
		err = errors.NewValidationError("Invalid limit").WithDetails("Limit must be between 1 and 100")
		return
	}

	offset = (page - 1) * limit
	return
}

// HandleError 统一的错误处理
func (bc *BaseController) HandleError(c *gin.Context, err error) {
	c.Error(err)
}
