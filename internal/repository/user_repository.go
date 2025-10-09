package repository

import (
	"context"
	"database/sql"
	"strings"

	"goMcp/internal/database"
	"goMcp/internal/database/generated/users"
	"goMcp/internal/dto"
	"goMcp/internal/errors"
	"goMcp/internal/utils"
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

	// 解析 FullName 为 FirstName 和 LastName
	var firstName, lastName *string
	if req.FullName != "" {
		// 简单的分割逻辑，可以根据需要改进
		parts := splitFullName(req.FullName)
		if len(parts) > 0 {
			firstName = &parts[0]
		}
		if len(parts) > 1 {
			lastName = &parts[1]
		}
	}

	// 创建用户参数
	params := users.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}
	
	if firstName != nil {
		params.FirstName = sql.NullString{String: *firstName, Valid: true}
	}
	if lastName != nil {
		params.LastName = sql.NullString{String: *lastName, Valid: true}
	}

	user, err := r.db.Users.CreateUser(ctx, params)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to create user", err)
	}

	return r.toUserResponse(user), nil
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	user, err := r.db.Users.GetUser(ctx, id)
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
	user, err := r.db.Users.GetUserByUsername(ctx, username)
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
	user, err := r.db.Users.GetUserByEmail(ctx, email)
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
	userList, err := r.db.Users.ListUsers(ctx, users.ListUsersParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to list users", err)
	}

	var responses []*dto.UserResponse
	for _, user := range userList {
		responses = append(responses, r.toUserResponse(user))
	}

	return responses, nil
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// 先获取当前用户信息
	currentUser, err := r.db.Users.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewUserNotFoundError()
		}
		return nil, errors.NewDatabaseError("Failed to get current user", err)
	}

	// 设置更新参数
	params := users.UpdateUserParams{
		ID: id,
	}

	if req.Email != nil {
		params.Email = *req.Email
	} else {
		params.Email = currentUser.Email
	}

	if req.FullName != nil {
		// 解析 FullName 为 FirstName 和 LastName
		parts := splitFullName(*req.FullName)
		if len(parts) > 0 {
			params.FirstName = sql.NullString{String: parts[0], Valid: true}
		} else {
			params.FirstName = sql.NullString{Valid: false}
		}
		if len(parts) > 1 {
			// 将除第一个部分外的所有部分作为LastName
			lastName := strings.Join(parts[1:], " ")
			params.LastName = sql.NullString{String: lastName, Valid: true}
		} else {
			params.LastName = sql.NullString{Valid: false}
		}
	} else {
		params.FirstName = currentUser.FirstName
		params.LastName = currentUser.LastName
	}

	if req.IsActive != nil {
		params.IsActive = sql.NullBool{Bool: *req.IsActive, Valid: true}
	} else {
		params.IsActive = currentUser.IsActive
	}

	user, err := r.db.Users.UpdateUser(ctx, params)
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
	// 首先检查用户是否存在
	_, err := r.db.Users.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NewUserNotFoundError()
		}
		return errors.NewDatabaseError("Failed to check user existence", err)
	}

	// 删除用户
	err = r.db.Users.DeleteUser(ctx, id)
	if err != nil {
		return errors.NewDatabaseError("Failed to delete user", err)
	}

	return nil
}

// ExistsByUsername 检查用户名是否存在
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	count, err := r.db.Users.CountUsersByUsername(ctx, username)
	if err != nil {
		return false, errors.NewDatabaseError("Failed to check username existence", err)
	}
	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := r.db.Users.CountUsersByEmail(ctx, email)
	if err != nil {
		return false, errors.NewDatabaseError("Failed to check email existence", err)
	}
	return count > 0, nil
}

// toUserResponse 将数据库用户模型转换为响应模型
func (r *userRepository) toUserResponse(user users.User) *dto.UserResponse {
	var fullName *string
	
	// 组合 FirstName 和 LastName 为 FullName
	if user.FirstName.Valid || user.LastName.Valid {
		var name string
		if user.FirstName.Valid {
			name = user.FirstName.String
		}
		if user.LastName.Valid {
			if name != "" {
				name += " " + user.LastName.String
			} else {
				name = user.LastName.String
			}
		}
		if name != "" {
			fullName = &name
		}
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

// splitFullName 将全名分割为名和姓
func splitFullName(fullName string) []string {
	if fullName == "" {
		return []string{}
	}
	
	// 使用strings.Fields来分割，它会自动处理多个空格
	return strings.Fields(fullName)
}
