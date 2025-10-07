package repository

import (
	"context"

	"admin/internal/database"
)

// repositoryManager 数据访问层管理器实现
type repositoryManager struct {
	db       *database.DB
	userRepo UserRepository
}

// NewRepositoryManager 创建数据访问层管理器
func NewRepositoryManager(db *database.DB) RepositoryManager {
	return &repositoryManager{
		db:       db,
		userRepo: NewUserRepository(db),
	}
}

// User 获取用户数据访问层
func (rm *repositoryManager) User() UserRepository {
	return rm.userRepo
}

// Close 关闭数据库连接
func (rm *repositoryManager) Close() error {
	return rm.db.Close()
}

// Ping 检查数据库连接
func (rm *repositoryManager) Ping(ctx context.Context) error {
	return rm.db.Ping(ctx)
}
