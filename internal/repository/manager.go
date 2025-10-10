package repository

import (
	"context"

	"go-springAi/internal/database"
)

// repositoryManager 数据访问层管理器实现
type repositoryManager struct {
	db         *database.DB
	userRepo   UserRepository
	apiKeyRepo APIKeyRepository
}

// NewRepositoryManager 创建数据访问层管理器
func NewRepositoryManager(db *database.DB) RepositoryManager {
	return &repositoryManager{
		db:         db,
		userRepo:   NewUserRepository(db),
		apiKeyRepo: NewAPIKeyRepository(db),
	}
}

// User 获取用户数据访问层
func (rm *repositoryManager) User() UserRepository {
	return rm.userRepo
}

// APIKey 获取API密钥数据访问层
func (rm *repositoryManager) APIKey() APIKeyRepository {
	return rm.apiKeyRepo
}

// Close 关闭数据库连接
func (rm *repositoryManager) Close() error {
	return rm.db.Close()
}

// Ping 检查数据库连接
func (rm *repositoryManager) Ping(ctx context.Context) error {
	return rm.db.Ping()
}
