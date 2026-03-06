package handler

import (
	"magickingdom-go/internal/dto"
	"magickingdom-go/internal/response"
	"magickingdom-go/internal/service"

	"github.com/gin-gonic/gin"
)

// SeatHandler 座位处理器
type SeatHandler struct {
	seatService          *service.SeatService
	seatOccupancyService *service.SeatOccupancyService
}

// NewSeatHandler 创建座位处理器实例
func NewSeatHandler(seatService *service.SeatService, seatOccupancyService *service.SeatOccupancyService) *SeatHandler {
	return &SeatHandler{
		seatService:          seatService,
		seatOccupancyService: seatOccupancyService,
	}
}

// GetSeatLayout 获取座位布局
// @Summary 获取座位布局
// @Description 获取酒吧的座位布局信息，包括DJ台、桌子、卡座、VIP和头等舱
// @Tags 座位
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=models.SeatLayout}
// @Router /api/v1/seats/layout [get]
func (h *SeatHandler) GetSeatLayout(c *gin.Context) {
	layout := h.seatService.GetSeatLayout()
	response.Success(c, layout)
}

// GetUserCurrentSeat 获取用户当前座位
// @Summary 获取用户当前座位
// @Description 获取当前登录用户的座位信息
// @Tags 座位
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=dto.UserSeatResponse}
// @Router /api/v1/seats/my-seat [get]
func (h *SeatHandler) GetUserCurrentSeat(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 401, "未授权")
		return
	}

	seatInfo, err := h.seatOccupancyService.GetUserCurrentSeat(userID.(uint))
	if err != nil {
		response.Error(c, 500, "获取座位信息失败: "+err.Error())
		return
	}

	response.Success(c, seatInfo)
}

// JoinSeat 用户入座
// @Summary 用户入座
// @Description 用户加入某个座位
// @Tags 座位
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.JoinSeatRequest true "入座请求"
// @Success 200 {object} response.Response
// @Router /api/v1/seats/join [post]
func (h *SeatHandler) JoinSeat(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 401, "未授权")
		return
	}

	var req dto.JoinSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := h.seatOccupancyService.JoinSeat(userID.(uint), &req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

// LeaveSeat 用户离座
// @Summary 用户离座
// @Description 用户离开当前座位
// @Tags 座位
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/v1/seats/leave [post]
func (h *SeatHandler) LeaveSeat(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 401, "未授权")
		return
	}

	if err := h.seatOccupancyService.LeaveSeat(userID.(uint)); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetSeatOccupancyInfo 获取座位占用信息
// @Summary 获取座位占用信息
// @Description 获取某个座位的入座人数和用户信息
// @Tags 座位
// @Accept json
// @Produce json
// @Param seat_id query string true "座位ID"
// @Success 200 {object} response.Response{data=models.SeatOccupancyInfo}
// @Router /api/v1/seats/occupancy [get]
func (h *SeatHandler) GetSeatOccupancyInfo(c *gin.Context) {
	seatID := c.Query("seat_id")
	if seatID == "" {
		response.Error(c, 400, "座位ID不能为空")
		return
	}

	info, err := h.seatOccupancyService.GetSeatOccupancyInfo(seatID)
	if err != nil {
		response.Error(c, 500, "获取座位信息失败: "+err.Error())
		return
	}

	response.Success(c, info)
}

