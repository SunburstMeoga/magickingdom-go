package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"magickingdom-go/internal/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	assert.NoError(t, err)

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		OpenID:   "test_open_id",
		Nickname: "Test User",
		Status:   1,
	}

	err := repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepository_FindByOpenID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	user := &models.User{
		OpenID:   "test_open_id",
		Nickname: "Test User",
		Status:   1,
	}
	err := repo.Create(user)
	assert.NoError(t, err)

	// 查找用户
	foundUser, err := repo.FindByOpenID("test_open_id")
	assert.NoError(t, err)
	assert.Equal(t, user.OpenID, foundUser.OpenID)
	assert.Equal(t, user.Nickname, foundUser.Nickname)
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	user := &models.User{
		OpenID:   "test_open_id",
		Nickname: "Test User",
		Status:   1,
	}
	err := repo.Create(user)
	assert.NoError(t, err)

	// 更新用户
	user.Nickname = "Updated User"
	err = repo.Update(user)
	assert.NoError(t, err)

	// 验证更新
	foundUser, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated User", foundUser.Nickname)
}

