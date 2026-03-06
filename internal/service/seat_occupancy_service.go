package service

import (
	"errors"
	"magickingdom-go/internal/dto"
	"magickingdom-go/internal/models"
	"magickingdom-go/internal/repository"

	"gorm.io/gorm"
)

// SeatOccupancyService 座位占用服务
type SeatOccupancyService struct {
	occupancyRepo *repository.SeatOccupancyRepository
}

// NewSeatOccupancyService 创建座位占用服务实例
func NewSeatOccupancyService(occupancyRepo *repository.SeatOccupancyRepository) *SeatOccupancyService {
	return &SeatOccupancyService{
		occupancyRepo: occupancyRepo,
	}
}

// GetUserCurrentSeat 获取用户当前座位
func (s *SeatOccupancyService) GetUserCurrentSeat(userID uint) (*dto.UserSeatResponse, error) {
	occupancy, err := s.occupancyRepo.GetUserCurrentSeat(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.UserSeatResponse{
				HasSeat: false,
			}, nil
		}
		return nil, err
	}

	return &dto.UserSeatResponse{
		HasSeat:  true,
		SeatID:   occupancy.SeatID,
		SeatType: occupancy.SeatType,
	}, nil
}

// JoinSeat 用户入座
func (s *SeatOccupancyService) JoinSeat(userID uint, req *dto.JoinSeatRequest) error {
	// 检查用户是否已经有座位
	currentSeat, err := s.occupancyRepo.GetUserCurrentSeat(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 如果用户已经在座位上
	if currentSeat != nil {
		// 如果是同一个座位，直接返回成功
		if currentSeat.SeatID == req.SeatID {
			return nil
		}
		// 如果是不同座位，返回错误
		return errors.New("用户已经在其他座位上，请先离座")
	}

	// 创建新的占用记录
	occupancy := &models.SeatOccupancy{
		UserID:   userID,
		SeatID:   req.SeatID,
		SeatType: req.SeatType,
		Status:   1,
	}

	return s.occupancyRepo.JoinSeat(occupancy)
}

// LeaveSeat 用户离座
func (s *SeatOccupancyService) LeaveSeat(userID uint) error {
	// 检查用户是否有座位
	_, err := s.occupancyRepo.GetUserCurrentSeat(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户当前没有座位")
		}
		return err
	}

	return s.occupancyRepo.LeaveSeat(userID)
}

// GetSeatOccupancyInfo 获取座位占用信息
func (s *SeatOccupancyService) GetSeatOccupancyInfo(seatID string) (*models.SeatOccupancyInfo, error) {
	occupancies, err := s.occupancyRepo.GetSeatOccupants(seatID)
	if err != nil {
		return nil, err
	}

	// 构建返回数据
	users := make([]models.SeatOccupancyUser, 0, len(occupancies))
	for _, occ := range occupancies {
		if occ.User != nil {
			users = append(users, models.SeatOccupancyUser{
				UserID:    occ.User.ID,
				Nickname:  occ.User.Nickname,
				AvatarURL: occ.User.AvatarURL,
				JoinedAt:  occ.CreatedAt,
			})
		}
	}

	// 获取座位类型（从第一条记录获取，如果没有记录则为空）
	seatType := ""
	if len(occupancies) > 0 {
		seatType = occupancies[0].SeatType
	}

	return &models.SeatOccupancyInfo{
		SeatID:      seatID,
		SeatType:    seatType,
		OccupiedNum: len(users),
		Users:       users,
	}, nil
}

