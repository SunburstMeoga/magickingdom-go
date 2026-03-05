package repository

import (
	"gorm.io/gorm"
	"magickingdom-go/internal/models"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByOpenID(openID string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}

// userRepository 用户仓储实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID 根据 ID 查找用户
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByOpenID 根据 OpenID 查找用户
func (r *userRepository) FindByOpenID(openID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("open_id = ?", openID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户（软删除）
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

