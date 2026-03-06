package dto

// JoinSeatRequest 入座请求
type JoinSeatRequest struct {
	SeatID   string `json:"seat_id" binding:"required"`   // 座位ID
	SeatType string `json:"seat_type" binding:"required"` // 座位类型
}

// LeaveSeatRequest 离座请求（可选，也可以不需要参数）
type LeaveSeatRequest struct {
	// 可以为空，从 JWT token 中获取用户信息
}

// UserSeatResponse 用户座位信息响应
type UserSeatResponse struct {
	HasSeat  bool   `json:"has_seat"`  // 是否有座位
	SeatID   string `json:"seat_id"`   // 座位ID
	SeatType string `json:"seat_type"` // 座位类型
}

