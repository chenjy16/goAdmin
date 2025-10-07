package services

import (
	"context"

	"admin/internal/dto"
	"admin/internal/errors"
	"admin/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to check username", err).WithStackTrace()
	}
	if exists {
		return nil, errors.NewConflictError("Username already exists")
	}

	// 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to check email", err).WithStackTrace()
	}
	if exists {
		return nil, errors.NewConflictError("Email already exists")
	}

	// 创建用户
	user, err := s.userRepo.Create(ctx, req)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to create user", err).WithStackTrace()
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		// 如果是数据库错误，包装为自定义错误
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.NewUserNotFoundError()
		}
		return nil, errors.NewDatabaseError("Failed to get user", err).WithStackTrace()
	}

	return user, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context, page, limit int64) ([]*dto.UserResponse, error) {
	params := repository.NewPaginationParams(page, limit)
	users, err := s.userRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.Update(ctx, id, req)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.NewUserNotFoundError()
		}
		return nil, errors.NewDatabaseError("Failed to update user", err).WithStackTrace()
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.NewUserNotFoundError()
		}
		return errors.NewDatabaseError("Failed to delete user", err).WithStackTrace()
	}

	return nil
}
