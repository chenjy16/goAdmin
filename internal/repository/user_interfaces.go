package repository

import (
	"context"

	"goMcp/internal/dto"
)

// PaginationParams 分页参数结构体
type PaginationParams struct {
	Page   int64 `json:"page"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

// NewPaginationParams 创建分页参数
func NewPaginationParams(page, limit int64) *PaginationParams {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	
	offset := (page - 1) * limit
	return &PaginationParams{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

// UserReader 用户读取接口
type UserReader interface {
	// 基础查询方法
	GetByID(ctx context.Context, id int64) (*dto.UserResponse, error)
	GetByUsername(ctx context.Context, username string) (*dto.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
	List(ctx context.Context, params *PaginationParams) ([]*dto.UserResponse, error)
}

// UserWriter 用户写入接口
type UserWriter interface {
	// 基础写入方法
	Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	Update(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id int64) error
}

// UserValidator 用户验证接口
type UserValidator interface {
	// 业务验证方法
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// UserRepository 用户数据访问接口（组合所有功能接口）
type UserRepository interface {
	UserReader
	UserWriter
	UserValidator
}

// RepositoryManager 数据访问层管理器接口
type RepositoryManager interface {
	User() UserRepository
	Close() error
	Ping(ctx context.Context) error
}