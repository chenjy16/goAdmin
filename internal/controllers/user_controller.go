package controllers

import (
	"net/http"

	"admin/internal/dto"
	"admin/internal/response"
	"admin/internal/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	*BaseController
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		BaseController: NewBaseController(),
		userService:    userService,
	}
}

// CreateUser 创建用户
func (uc *UserController) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	// 使用基础控制器的统一绑定和验证方法
	if err := uc.BindAndValidate(c, &req); err != nil {
		return
	}

	user, err := uc.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, "User created successfully", user)
}

// GetUser 获取单个用户
func (uc *UserController) GetUser(c *gin.Context) {
	// 使用基础控制器的ID解析方法
	id, err := uc.ParseIDParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	user, err := uc.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

// ListUsers 获取用户列表
func (uc *UserController) ListUsers(c *gin.Context) {
	// 使用基础控制器的分页参数解析方法
	page, limit, _, err := uc.ParsePaginationParams(c)
	if err != nil {
		c.Error(err)
		return
	}

	users, err := uc.userService.ListUsers(c.Request.Context(), page, limit)
	if err != nil {
		c.Error(err)
		return
	}

	// 构建分页响应
	paginationData := map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": len(users), // 注意：这里暂时使用当前页的用户数，实际应该从服务层获取总数
		},
	}

	response.Success(c, http.StatusOK, "Users retrieved successfully", paginationData)
}

// UpdateUser 更新用户
func (uc *UserController) UpdateUser(c *gin.Context) {
	// 使用基础控制器的ID解析方法
	id, err := uc.ParseIDParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdateUserRequest

	// 使用基础控制器的统一绑定和验证方法
	if err := uc.BindAndValidate(c, &req); err != nil {
		return
	}

	user, err := uc.userService.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "User updated successfully", user)
}

// DeleteUser 删除用户
func (uc *UserController) DeleteUser(c *gin.Context) {
	// 使用基础控制器的ID解析方法
	id, err := uc.ParseIDParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	err = uc.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}
