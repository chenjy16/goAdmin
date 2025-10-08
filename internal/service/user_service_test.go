package service

import (
	"testing"

	"admin/internal/dto"
	"admin/internal/errors"
	"admin/internal/mocks"
	"admin/internal/testutil"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_Create(t *testing.T) {
	testutil.SetupGinTest()

	tests := []struct {
		name    string
		request dto.CreateUserRequest
		setup   func(*mocks.MockUserRepository)
		wantErr bool
		errType *errors.AppError
	}{
		{
			name: "Success",
			request: dto.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				FullName: "Test User",
			},
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(&dto.UserResponse{
						ID:       1,
						Username: "testuser",
						Email:    "test@example.com",
						FullName: stringPtr("Test User"),
					}, nil)
			},
			wantErr: false,
		},
		{
			name: "Error_DuplicateUsername",
			request: dto.CreateUserRequest{
				Username: "existinguser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.NewUsernameExistsError())
			},
			wantErr: true,
			errType: errors.NewUsernameExistsError(),
		},
		{
			name: "Error_DuplicateEmail",
			request: dto.CreateUserRequest{
				Username: "testuser",
				Email:    "existing@example.com",
				Password: "password123",
			},
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.NewEmailExistsError())
			},
			wantErr: true,
			errType: errors.NewEmailExistsError(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 创建mock repository
			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			mockRepoManager := mocks.NewMockRepositoryManager(ctrl)
			mockRepoManager.EXPECT().User().Return(mockUserRepo).AnyTimes()

			// 设置mock期望
			tt.setup(mockUserRepo)

			// 创建服务
			userService := NewUserService(mockRepoManager)

			// 执行测试
			ctx := testutil.TestContext()
			result, err := userService.Create(ctx, tt.request)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errType != nil {
					assert.IsType(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Username, result.Username)
				assert.Equal(t, tt.request.Email, result.Email)
			}
		})
	}
}

func TestUserService_GetByID(t *testing.T) {
	testutil.SetupGinTest()

	tests := []struct {
		name    string
		userID  int64
		setup   func(*mocks.MockUserRepository)
		wantErr bool
		errType *errors.AppError
	}{
		{
			name:   "Success",
			userID: 1,
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().GetByID(gomock.Any(), int64(1)).
					Return(&dto.UserResponse{
						ID:       1,
						Username: "testuser",
						Email:    "test@example.com",
					}, nil)
			},
			wantErr: false,
		},
		{
			name:   "Error_UserNotFound",
			userID: 999,
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().GetByID(gomock.Any(), int64(999)).
					Return(nil, errors.NewUserNotFoundError())
			},
			wantErr: true,
			errType: errors.NewUserNotFoundError(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 创建mock repository
			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			mockRepoManager := mocks.NewMockRepositoryManager(ctrl)
			mockRepoManager.EXPECT().User().Return(mockUserRepo).AnyTimes()

			// 设置mock期望
			tt.setup(mockUserRepo)

			// 创建服务
			userService := NewUserService(mockRepoManager)

			// 执行测试
			ctx := testutil.TestContext()
			result, err := userService.GetByID(ctx, tt.userID)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errType != nil {
					assert.IsType(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.userID, result.ID)
			}
		})
	}
}

func TestUserService_GetByUsername(t *testing.T) {
	testutil.SetupGinTest()

	tests := []struct {
		name     string
		username string
		setup    func(*mocks.MockUserRepository)
		wantErr  bool
		errType  *errors.AppError
	}{
		{
			name:     "Success",
			username: "testuser",
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().GetByUsername(gomock.Any(), "testuser").
					Return(&dto.UserResponse{
						ID:       1,
						Username: "testuser",
						Email:    "test@example.com",
					}, nil)
			},
			wantErr: false,
		},
		{
			name:     "Error_UserNotFound",
			username: "nonexistent",
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().GetByUsername(gomock.Any(), "nonexistent").
					Return(nil, errors.NewUserNotFoundError())
			},
			wantErr: true,
			errType: errors.NewUserNotFoundError(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 创建mock repository
			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			mockRepoManager := mocks.NewMockRepositoryManager(ctrl)
			mockRepoManager.EXPECT().User().Return(mockUserRepo).AnyTimes()

			// 设置mock期望
			tt.setup(mockUserRepo)

			// 创建服务
			userService := NewUserService(mockRepoManager)

			// 执行测试
			ctx := testutil.TestContext()
			result, err := userService.GetByUsername(ctx, tt.username)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errType != nil {
					assert.IsType(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.username, result.Username)
			}
		})
	}
}

func TestUserService_List(t *testing.T) {
	testutil.SetupGinTest()

	tests := []struct {
		name    string
		page    int64
		limit   int64
		setup   func(*mocks.MockUserRepository)
		wantErr bool
	}{
		{
			name:  "Success",
			page:  1,
			limit: 10,
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).
					Return([]*dto.UserResponse{
						{ID: 1, Username: "user1", Email: "user1@example.com"},
						{ID: 2, Username: "user2", Email: "user2@example.com"},
					}, nil)
			},
			wantErr: false,
		},
		{
			name:  "Success_EmptyList",
			page:  1,
			limit: 10,
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).
					Return([]*dto.UserResponse{}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 创建mock repository
			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			mockRepoManager := mocks.NewMockRepositoryManager(ctrl)
			mockRepoManager.EXPECT().User().Return(mockUserRepo).AnyTimes()

			// 设置mock期望
			tt.setup(mockUserRepo)

			// 创建服务
			userService := NewUserService(mockRepoManager)

			// 执行测试
			ctx := testutil.TestContext()
			result, err := userService.List(ctx, tt.page, tt.limit)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestUserService_Update(t *testing.T) {
	testutil.SetupGinTest()

	tests := []struct {
		name    string
		userID  int64
		request dto.UpdateUserRequest
		setup   func(*mocks.MockUserRepository)
		wantErr bool
		errType *errors.AppError
	}{
		{
			name:   "Success",
			userID: 1,
			request: dto.UpdateUserRequest{
				Email:    stringPtr("updated@example.com"),
				FullName: stringPtr("Updated User"),
			},
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().Update(gomock.Any(), int64(1), gomock.Any()).
					Return(&dto.UserResponse{
						ID:       1,
						Username: "testuser",
						Email:    "updated@example.com",
						FullName: stringPtr("Updated User"),
					}, nil)
			},
			wantErr: false,
		},
		{
			name:   "Error_UserNotFound",
			userID: 999,
			request: dto.UpdateUserRequest{
				Email: stringPtr("updated@example.com"),
			},
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().Update(gomock.Any(), int64(999), gomock.Any()).
					Return(nil, errors.NewUserNotFoundError())
			},
			wantErr: true,
			errType: errors.NewUserNotFoundError(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 创建mock repository
			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			mockRepoManager := mocks.NewMockRepositoryManager(ctrl)
			mockRepoManager.EXPECT().User().Return(mockUserRepo).AnyTimes()

			// 设置mock期望
			tt.setup(mockUserRepo)

			// 创建服务
			userService := NewUserService(mockRepoManager)

			// 执行测试
			ctx := testutil.TestContext()
			result, err := userService.Update(ctx, tt.userID, tt.request)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errType != nil {
					assert.IsType(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.userID, result.ID)
			}
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	testutil.SetupGinTest()

	tests := []struct {
		name    string
		userID  int64
		setup   func(*mocks.MockUserRepository)
		wantErr bool
		errType *errors.AppError
	}{
		{
			name:   "Success",
			userID: 1,
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().Delete(gomock.Any(), int64(1)).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Error_UserNotFound",
			userID: 999,
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().Delete(gomock.Any(), int64(999)).
					Return(errors.NewUserNotFoundError())
			},
			wantErr: true,
			errType: errors.NewUserNotFoundError(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 创建mock repository
			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			mockRepoManager := mocks.NewMockRepositoryManager(ctrl)
			mockRepoManager.EXPECT().User().Return(mockUserRepo).AnyTimes()

			// 设置mock期望
			tt.setup(mockUserRepo)

			// 创建服务
			userService := NewUserService(mockRepoManager)

			// 执行测试
			ctx := testutil.TestContext()
			err := userService.Delete(ctx, tt.userID)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.IsType(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_ExistsByUsername(t *testing.T) {
	testutil.SetupGinTest()

	tests := []struct {
		name     string
		username string
		setup    func(*mocks.MockUserRepository)
		expected bool
		wantErr  bool
	}{
		{
			name:     "Success_Exists",
			username: "existinguser",
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().ExistsByUsername(gomock.Any(), "existinguser").
					Return(true, nil)
			},
			expected: true,
			wantErr:  false,
		},
		{
			name:     "Success_NotExists",
			username: "nonexistentuser",
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().ExistsByUsername(gomock.Any(), "nonexistentuser").
					Return(false, nil)
			},
			expected: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 创建mock repository
			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			mockRepoManager := mocks.NewMockRepositoryManager(ctrl)
			mockRepoManager.EXPECT().User().Return(mockUserRepo).AnyTimes()

			// 设置mock期望
			tt.setup(mockUserRepo)

			// 创建服务
			userService := NewUserService(mockRepoManager)

			// 执行测试
			ctx := testutil.TestContext()
			result, err := userService.ExistsByUsername(ctx, tt.username)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestUserService_ExistsByEmail(t *testing.T) {
	testutil.SetupGinTest()

	tests := []struct {
		name     string
		email    string
		setup    func(*mocks.MockUserRepository)
		expected bool
		wantErr  bool
	}{
		{
			name:  "Success_Exists",
			email: "existing@example.com",
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().ExistsByEmail(gomock.Any(), "existing@example.com").
					Return(true, nil)
			},
			expected: true,
			wantErr:  false,
		},
		{
			name:  "Success_NotExists",
			email: "nonexistent@example.com",
			setup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().ExistsByEmail(gomock.Any(), "nonexistent@example.com").
					Return(false, nil)
			},
			expected: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 创建mock repository
			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			mockRepoManager := mocks.NewMockRepositoryManager(ctrl)
			mockRepoManager.EXPECT().User().Return(mockUserRepo).AnyTimes()

			// 设置mock期望
			tt.setup(mockUserRepo)

			// 创建服务
			userService := NewUserService(mockRepoManager)

			// 执行测试
			ctx := testutil.TestContext()
			result, err := userService.ExistsByEmail(ctx, tt.email)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// stringPtr 辅助函数，返回字符串指针
func stringPtr(s string) *string {
	return &s
}