package handler

import (
	"magickingdom-go/internal/response"
	"magickingdom-go/internal/service"

	"github.com/gin-gonic/gin"
)

// SeatHandler 座位处理器
type SeatHandler struct {
	seatService *service.SeatService
}

// NewSeatHandler 创建座位处理器实例
func NewSeatHandler(seatService *service.SeatService) *SeatHandler {
	return &SeatHandler{
		seatService: seatService,
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

