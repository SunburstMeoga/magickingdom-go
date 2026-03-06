package repository

import (
	"magickingdom-go/internal/models"

	"gorm.io/gorm"
)

// SeatOccupancyRepository 座位占用仓储
type SeatOccupancyRepository struct {
	db *gorm.DB
}

// NewSeatOccupancyRepository 创建座位占用仓储实例
func NewSeatOccupancyRepository(db *gorm.DB) *SeatOccupancyRepository {
	return &SeatOccupancyRepository{db: db}
}

// GetUserCurrentSeat 获取用户当前占用的座位
func (r *SeatOccupancyRepository) GetUserCurrentSeat(userID uint) (*models.SeatOccupancy, error) {
	var occupancy models.SeatOccupancy
	err := r.db.Where("user_id = ? AND status = ?", userID, 1).First(&occupancy).Error
	if err != nil {
		return nil, err
	}
	return &occupancy, nil
}

// JoinSeat 用户入座
func (r *SeatOccupancyRepository) JoinSeat(occupancy *models.SeatOccupancy) error {
	return r.db.Create(occupancy).Error
}

// LeaveSeat 用户离座（软删除或更新状态）
func (r *SeatOccupancyRepository) LeaveSeat(userID uint) error {
	return r.db.Model(&models.SeatOccupancy{}).
		Where("user_id = ? AND status = ?", userID, 1).
		Updates(map[string]interface{}{"status": 0}).Error
}

// GetSeatOccupants 获取某个座位的所有占用者
func (r *SeatOccupancyRepository) GetSeatOccupants(seatID string) ([]models.SeatOccupancy, error) {
	var occupancies []models.SeatOccupancy
	err := r.db.Preload("User").
		Where("seat_id = ? AND status = ?", seatID, 1).
		Order("created_at ASC").
		Find(&occupancies).Error
	return occupancies, err
}

// CountSeatOccupants 统计某个座位的占用人数
func (r *SeatOccupancyRepository) CountSeatOccupants(seatID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.SeatOccupancy{}).
		Where("seat_id = ? AND status = ?", seatID, 1).
		Count(&count).Error
	return count, err
}

