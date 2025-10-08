package repository

import (
	"context"
	"database/sql"
	"testing"

	"admin/internal/database"
	"admin/internal/database/generated/users"
	"admin/internal/dto"
	"admin/internal/errors"
	"admin/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_GetByID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := testutil.TestContext()

	// 创建测试用户
	testUser := createTestUser(t, db)

	t.Run("Success", func(t *testing.T) {
		user, err := repo.GetByID(ctx, testUser.ID)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Username, user.Username)
		assert.Equal(t, testUser.Email, user.Email)
		assert.Equal(t, testUser.IsActive.Bool, user.IsActive)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		user, err := repo.GetByID(ctx, 99999)
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.IsType(t, &errors.AppError{}, err)
	})
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := testutil.TestContext()

	// 创建测试用户
	testUser := createTestUser(t, db)

	t.Run("Success", func(t *testing.T) {
		user, err := repo.GetByUsername(ctx, testUser.Username)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Username, user.Username)
		assert.Equal(t, testUser.Email, user.Email)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		user, err := repo.GetByUsername(ctx, "nonexistent")
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.IsType(t, &errors.AppError{}, err)
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := testutil.TestContext()

	// 创建测试用户
	testUser := createTestUser(t, db)

	t.Run("Success", func(t *testing.T) {
		user, err := repo.GetByEmail(ctx, testUser.Email)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Username, user.Username)
		assert.Equal(t, testUser.Email, user.Email)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		user, err := repo.GetByEmail(ctx, "nonexistent@example.com")
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.IsType(t, &errors.AppError{}, err)
	})
}

func TestUserRepository_List(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := testutil.TestContext()

	// 创建多个测试用户
	user1 := createTestUser(t, db)
	user2 := createTestUserWithData(t, db, "testuser2", "test2@example.com")
	user3 := createTestUserWithData(t, db, "testuser3", "test3@example.com")

	t.Run("Success_WithPagination", func(t *testing.T) {
		params := NewPaginationParams(1, 2)
		users, err := repo.List(ctx, params)
		require.NoError(t, err)
		assert.Len(t, users, 2)
		
		// 验证分页参数
		assert.Equal(t, int64(1), params.Page)
		assert.Equal(t, int64(2), params.Limit)
		assert.Equal(t, int64(0), params.Offset)
	})

	t.Run("Success_SecondPage", func(t *testing.T) {
		params := NewPaginationParams(2, 2)
		users, err := repo.List(ctx, params)
		require.NoError(t, err)
		assert.Len(t, users, 1) // 第二页只有一个用户
		
		// 验证分页参数
		assert.Equal(t, int64(2), params.Page)
		assert.Equal(t, int64(2), params.Limit)
		assert.Equal(t, int64(2), params.Offset)
	})

	t.Run("Success_EmptyResult", func(t *testing.T) {
		params := NewPaginationParams(10, 10)
		users, err := repo.List(ctx, params)
		require.NoError(t, err)
		assert.Len(t, users, 0)
	})

	// 清理测试数据
	_ = user1
	_ = user2
	_ = user3
}

func TestUserRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := testutil.TestContext()

	t.Run("Success", func(t *testing.T) {
		req := dto.CreateUserRequest{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "password123",
			FullName: "New User",
		}

		user, err := repo.Create(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotZero(t, user.ID)
		assert.Equal(t, req.Username, user.Username)
		assert.Equal(t, req.Email, user.Email)
		assert.Equal(t, req.FullName, *user.FullName)
		assert.True(t, user.IsActive)
		assert.NotZero(t, user.CreatedAt)
		assert.NotZero(t, user.UpdatedAt)
	})

	t.Run("Success_WithoutFullName", func(t *testing.T) {
		req := dto.CreateUserRequest{
			Username: "newuser2",
			Email:    "newuser2@example.com",
			Password: "password123",
		}

		user, err := repo.Create(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Nil(t, user.FullName)
	})

	t.Run("Error_DuplicateUsername", func(t *testing.T) {
		// 先创建一个用户
		req1 := dto.CreateUserRequest{
			Username: "duplicate",
			Email:    "duplicate1@example.com",
			Password: "password123",
		}
		_, err := repo.Create(ctx, req1)
		require.NoError(t, err)

		// 尝试创建相同用户名的用户
		req2 := dto.CreateUserRequest{
			Username: "duplicate",
			Email:    "duplicate2@example.com",
			Password: "password123",
		}
		user, err := repo.Create(ctx, req2)
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.IsType(t, &errors.AppError{}, err)
	})

	t.Run("Error_DuplicateEmail", func(t *testing.T) {
		// 先创建一个用户
		req1 := dto.CreateUserRequest{
			Username: "user1",
			Email:    "duplicate@example.com",
			Password: "password123",
		}
		_, err := repo.Create(ctx, req1)
		require.NoError(t, err)

		// 尝试创建相同邮箱的用户
		req2 := dto.CreateUserRequest{
			Username: "user2",
			Email:    "duplicate@example.com",
			Password: "password123",
		}
		user, err := repo.Create(ctx, req2)
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.IsType(t, &errors.AppError{}, err)
	})
}

func TestUserRepository_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := testutil.TestContext()

	// 创建测试用户
	testUser := createTestUser(t, db)

	t.Run("Success_UpdateEmail", func(t *testing.T) {
		newEmail := "updated@example.com"
		req := dto.UpdateUserRequest{
			Email: &newEmail,
		}

		user, err := repo.Update(ctx, testUser.ID, req)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, newEmail, user.Email)
		assert.Equal(t, testUser.Username, user.Username) // 其他字段不变
	})

	t.Run("Success_UpdateFullName", func(t *testing.T) {
		newFullName := "Updated Full Name"
		req := dto.UpdateUserRequest{
			FullName: &newFullName,
		}

		user, err := repo.Update(ctx, testUser.ID, req)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, newFullName, *user.FullName)
	})

	t.Run("Success_UpdateIsActive", func(t *testing.T) {
		isActive := false
		req := dto.UpdateUserRequest{
			IsActive: &isActive,
		}

		user, err := repo.Update(ctx, testUser.ID, req)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.False(t, user.IsActive)
	})

	t.Run("Success_UpdateMultipleFields", func(t *testing.T) {
		newEmail := "multi@example.com"
		newFullName := "Multi Update"
		isActive := false
		req := dto.UpdateUserRequest{
			Email:    &newEmail,
			FullName: &newFullName,
			IsActive: &isActive,
		}

		user, err := repo.Update(ctx, testUser.ID, req)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, newEmail, user.Email)
		assert.Equal(t, newFullName, *user.FullName)
		assert.False(t, user.IsActive)
	})

	t.Run("Error_UserNotFound", func(t *testing.T) {
		newEmail := "notfound@example.com"
		req := dto.UpdateUserRequest{
			Email: &newEmail,
		}

		user, err := repo.Update(ctx, 99999, req)
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.IsType(t, &errors.AppError{}, err)
	})
}

func TestUserRepository_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := testutil.TestContext()

	t.Run("Success", func(t *testing.T) {
		// 创建测试用户
		testUser := createTestUser(t, db)

		err := repo.Delete(ctx, testUser.ID)
		require.NoError(t, err)

		// 验证用户已被删除
		user, err := repo.GetByID(ctx, testUser.ID)
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.IsType(t, &errors.AppError{}, err)
	})

	t.Run("Error_UserNotFound", func(t *testing.T) {
		err := repo.Delete(ctx, 99999)
		assert.Error(t, err)
		assert.IsType(t, &errors.AppError{}, err)
	})
}

func TestUserRepository_ExistsByUsername(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := testutil.TestContext()

	// 创建测试用户
	testUser := createTestUser(t, db)

	t.Run("Success_Exists", func(t *testing.T) {
		exists, err := repo.ExistsByUsername(ctx, testUser.Username)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Success_NotExists", func(t *testing.T) {
		exists, err := repo.ExistsByUsername(ctx, "nonexistent")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := testutil.TestContext()

	// 创建测试用户
	testUser := createTestUser(t, db)

	t.Run("Success_Exists", func(t *testing.T) {
		exists, err := repo.ExistsByEmail(ctx, testUser.Email)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Success_NotExists", func(t *testing.T) {
		exists, err := repo.ExistsByEmail(ctx, "nonexistent@example.com")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestPaginationParams(t *testing.T) {
	t.Run("ValidParams", func(t *testing.T) {
		params := NewPaginationParams(2, 10)
		assert.Equal(t, int64(2), params.Page)
		assert.Equal(t, int64(10), params.Limit)
		assert.Equal(t, int64(10), params.Offset) // (2-1) * 10
	})

	t.Run("InvalidPage_ShouldDefault", func(t *testing.T) {
		params := NewPaginationParams(0, 10)
		assert.Equal(t, int64(1), params.Page)
		assert.Equal(t, int64(0), params.Offset)
	})

	t.Run("InvalidLimit_ShouldDefault", func(t *testing.T) {
		params := NewPaginationParams(1, 0)
		assert.Equal(t, int64(10), params.Limit)
	})

	t.Run("LimitTooLarge_ShouldCap", func(t *testing.T) {
		params := NewPaginationParams(1, 200)
		assert.Equal(t, int64(100), params.Limit)
	})
}

// 辅助函数：创建测试用户
func createTestUser(t *testing.T, db *database.DB) users.User {
	return createTestUserWithData(t, db, "testuser", "test@example.com")
}

func createTestUserWithData(t *testing.T, db *database.DB, username, email string) users.User {
	t.Helper()

	user, err := db.Users.CreateUser(context.Background(), users.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: "hashed_password",
		FirstName:    sql.NullString{String: "Test", Valid: true},
		LastName:     sql.NullString{String: "User", Valid: true},
	})
	require.NoError(t, err)
	return user
}