package repository

import (
	"testing"

	"admin/internal/dto"
	"admin/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepositoryManager_NewRepositoryManager(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	manager := NewRepositoryManager(db)
	assert.NotNil(t, manager)
}

func TestRepositoryManager_User(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	manager := NewRepositoryManager(db)
	userRepo := manager.User()
	
	assert.NotNil(t, userRepo)
	assert.Implements(t, (*UserRepository)(nil), userRepo)
}

func TestRepositoryManager_Ping(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	manager := NewRepositoryManager(db)
	ctx := testutil.TestContext()

	t.Run("Success", func(t *testing.T) {
		err := manager.Ping(ctx)
		require.NoError(t, err)
	})
}

func TestRepositoryManager_Close(t *testing.T) {
	db := testutil.SetupTestDB(t)
	
	manager := NewRepositoryManager(db)

	t.Run("Success", func(t *testing.T) {
		err := manager.Close()
		require.NoError(t, err)
	})
}

func TestRepositoryManager_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	manager := NewRepositoryManager(db)
	ctx := testutil.TestContext()

	t.Run("UserRepository_Integration", func(t *testing.T) {
		userRepo := manager.User()
		
		// 测试创建用户
		createReq := dto.CreateUserRequest{
			Username: "integration_test",
			Email:    "integration@example.com",
			Password: "password123",
			FullName: "Integration Test User",
		}
		
		user, err := userRepo.Create(ctx, createReq)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, createReq.Username, user.Username)
		assert.Equal(t, createReq.Email, user.Email)
		
		// 测试获取用户
		retrievedUser, err := userRepo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, retrievedUser.ID)
		assert.Equal(t, user.Username, retrievedUser.Username)
		
		// 测试用户存在性检查
		exists, err := userRepo.ExistsByUsername(ctx, user.Username)
		require.NoError(t, err)
		assert.True(t, exists)
		
		exists, err = userRepo.ExistsByEmail(ctx, user.Email)
		require.NoError(t, err)
		assert.True(t, exists)
		
		// 测试更新用户
		newEmail := "updated_integration@example.com"
		updateReq := dto.UpdateUserRequest{
			Email: &newEmail,
		}
		
		updatedUser, err := userRepo.Update(ctx, user.ID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, newEmail, updatedUser.Email)
		
		// 测试删除用户
		err = userRepo.Delete(ctx, user.ID)
		require.NoError(t, err)
		
		// 验证用户已被删除
		deletedUser, err := userRepo.GetByID(ctx, user.ID)
		assert.Nil(t, deletedUser)
		assert.Error(t, err)
	})
}