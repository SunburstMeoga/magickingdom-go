package dto

// WechatLoginRequest 微信登录请求
type WechatLoginRequest struct {
	Code string `json:"code" binding:"required"` // 微信登录凭证
}

// WechatLoginResponse 微信登录响应
type WechatLoginResponse struct {
	Token string      `json:"token"`
	User  UserInfoDTO `json:"user"`
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Gender    *int   `json:"gender"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Language  string `json:"language"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

// UserInfoDTO 用户信息 DTO
type UserInfoDTO struct {
	ID        uint   `json:"id"`
	OpenID    string `json:"open_id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Gender    int    `json:"gender"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Language  string `json:"language"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	Status    int    `json:"status"`
}

// TestTokenRequest 测试 Token 请求
type TestTokenRequest struct {
	UserID uint `json:"user_id" binding:"required"` // 用户ID
}

