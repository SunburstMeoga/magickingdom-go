package models

import (
	"time"

	"gorm.io/gorm"
)

// SeatOccupancy 座位占用记录
type SeatOccupancy struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 用户信息
	UserID uint   `gorm:"not null;index" json:"user_id"`
	User   *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// 座位信息
	SeatID   string `gorm:"size:50;not null;index" json:"seat_id"`     // 座位ID，如 C09+, V19+, T01 等
	SeatType string `gorm:"size:20;not null;index" json:"seat_type"`   // 座位类型：card, vip, table, first-class
	
	// 状态
	Status int `gorm:"default:1;index" json:"status"` // 1-占用中 0-已离座
}

// TableName 指定表名
func (SeatOccupancy) TableName() string {
	return "seat_occupancies"
}

// SeatOccupancyInfo 座位占用信息（用于返回给前端）
type SeatOccupancyInfo struct {
	SeatID      string                `json:"seat_id"`
	SeatType    string                `json:"seat_type"`
	OccupiedNum int                   `json:"occupied_num"` // 当前入座人数
	Users       []SeatOccupancyUser   `json:"users"`        // 入座用户列表
}

// SeatOccupancyUser 入座用户信息
type SeatOccupancyUser struct {
	UserID    uint      `json:"user_id"`
	Nickname  string    `json:"nickname"`
	AvatarURL string    `json:"avatar_url"`
	JoinedAt  time.Time `json:"joined_at"` // 入座时间
}

