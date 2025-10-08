package service

import (
	"context"

	"admin/internal/dto"
	"admin/internal/repository"
)

// UserService 用户服务接口
type UserService interface {
	// Create 创建用户
	Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id int64) (*dto.UserResponse, error)
	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*dto.UserResponse, error)
	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
	// List 获取用户列表
	List(ctx context.Context, page, limit int64) ([]*dto.UserResponse, error)
	// Update 更新用户
	Update(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	// Delete 删除用户
	Delete(ctx context.Context, id int64) error
	// ExistsByUsername 检查用户名是否存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(repoManager repository.RepositoryManager) UserService {
	return &userService{
		userRepo: repoManager.User(),
	}
}

// Create 创建用户
func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	return s.userRepo.Create(ctx, req)
}

// GetByID 根据ID获取用户
func (s *userService) GetByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetByUsername 根据用户名获取用户
func (s *userService) GetByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

// GetByEmail 根据邮箱获取用户
func (s *userService) GetByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// List 获取用户列表
func (s *userService) List(ctx context.Context, page, limit int64) ([]*dto.UserResponse, error) {
	params := repository.NewPaginationParams(page, limit)
	return s.userRepo.List(ctx, params)
}

// Update 更新用户
func (s *userService) Update(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	return s.userRepo.Update(ctx, id, req)
}

// Delete 删除用户
func (s *userService) Delete(ctx context.Context, id int64) error {
	return s.userRepo.Delete(ctx, id)
}

// ExistsByUsername 检查用户名是否存在
func (s *userService) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return s.userRepo.ExistsByUsername(ctx, username)
}

// ExistsByEmail 检查邮箱是否存在
func (s *userService) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return s.userRepo.ExistsByEmail(ctx, email)
}