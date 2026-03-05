package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 微信相关字段
	OpenID     string `gorm:"uniqueIndex;size:100;not null" json:"open_id"`
	UnionID    string `gorm:"index;size:100" json:"union_id,omitempty"`
	SessionKey string `gorm:"size:100" json:"-"` // 不返回给前端

	// 用户信息
	Nickname  string `gorm:"size:100" json:"nickname"`
	AvatarURL string `gorm:"size:500" json:"avatar_url"`
	Gender    int    `gorm:"default:0" json:"gender"` // 0-未知 1-男 2-女
	Country   string `gorm:"size:50" json:"country"`
	Province  string `gorm:"size:50" json:"province"`
	City      string `gorm:"size:50" json:"city"`
	Language  string `gorm:"size:20" json:"language"`

	// 业务字段
	Phone  string `gorm:"index;size:20" json:"phone,omitempty"`
	Email  string `gorm:"size:100" json:"email,omitempty"`
	Status int    `gorm:"default:1" json:"status"` // 1-正常 0-禁用
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

