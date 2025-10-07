package repository

import (
	"context"
	"database/sql"

	"admin/internal/database"
	"admin/internal/dto"
	"admin/internal/errors"
	"admin/internal/utils"
)

// userRepository 用户数据访问层实现
type userRepository struct {
	db *database.DB
}

// NewUserRepository 创建用户数据访问层实例
func NewUserRepository(db *database.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewInternalError("Failed to hash password").WithCause(err)
	}

	// 创建用户
	var fullName sql.NullString
	if req.FullName != "" {
		fullName = sql.NullString{String: req.FullName, Valid: true}
	}

	user, err := r.db.CreateUser(ctx, database.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     fullName,
	})
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to create user", err)
	}

	return r.toUserResponse(user), nil
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	user, err := r.db.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewUserNotFoundError()
		}
		return nil, errors.NewDatabaseError("Failed to get user by ID", err)
	}

	return r.toUserResponse(user), nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
	user, err := r.db.GetUserByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewUserNotFoundError()
		}
		return nil, errors.NewDatabaseError("Failed to get user by username", err)
	}

	return r.toUserResponse(user), nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	user, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewUserNotFoundError()
		}
		return nil, errors.NewDatabaseError("Failed to get user by email", err)
	}

	return r.toUserResponse(user), nil
}

// List 获取用户列表
func (r *userRepository) List(ctx context.Context, params *PaginationParams) ([]*dto.UserResponse, error) {
	users, err := r.db.ListUsers(ctx, database.ListUsersParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to list users", err)
	}

	var responses []*dto.UserResponse
	for _, user := range users {
		responses = append(responses, r.toUserResponse(user))
	}

	return responses, nil
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	var email string
	var fullName sql.NullString
	var isActive sql.NullBool

	// 先获取当前用户信息
	currentUser, err := r.db.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewUserNotFoundError()
		}
		return nil, errors.NewDatabaseError("Failed to get current user", err)
	}

	// 设置更新参数
	if req.Email != nil {
		email = *req.Email
	} else {
		email = currentUser.Email
	}

	if req.FullName != nil {
		fullName = sql.NullString{String: *req.FullName, Valid: true}
	} else {
		fullName = currentUser.FullName
	}

	if req.IsActive != nil {
		isActive = sql.NullBool{Bool: *req.IsActive, Valid: true}
	} else {
		isActive = currentUser.IsActive
	}

	user, err := r.db.UpdateUser(ctx, database.UpdateUserParams{
		ID:       id,
		Email:    email,
		FullName: fullName,
		IsActive: isActive,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewUserNotFoundError()
		}
		return nil, errors.NewDatabaseError("Failed to update user", err)
	}

	return r.toUserResponse(user), nil
}

// Delete 删除用户
func (r *userRepository) Delete(ctx context.Context, id int64) error {
	err := r.db.DeleteUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NewUserNotFoundError()
		}
		return errors.NewDatabaseError("Failed to delete user", err)
	}

	return nil
}

// ExistsByUsername 检查用户名是否存在
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	_, err := r.db.GetUserByUsername(ctx, username)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, errors.NewDatabaseError("Failed to check username existence", err)
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	_, err := r.db.GetUserByEmail(ctx, email)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, errors.NewDatabaseError("Failed to check email existence", err)
}



// toUserResponse 将数据库用户模型转换为响应模型
func (r *userRepository) toUserResponse(user database.User) *dto.UserResponse {
	var fullName *string
	if user.FullName.Valid {
		fullName = &user.FullName.String
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  fullName,
		IsActive:  user.IsActive.Bool,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}
}
