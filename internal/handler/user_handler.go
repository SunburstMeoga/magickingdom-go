package handler

import (
	"github.com/gin-gonic/gin"
	"magickingdom-go/internal/dto"
	"magickingdom-go/internal/middleware"
	"magickingdom-go/internal/response"
	"magickingdom-go/internal/service"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// WechatLogin 微信登录
// @Summary 微信小程序登录
// @Description 使用微信小程序的 code 进行登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body dto.WechatLoginRequest true "登录请求"
// @Success 200 {object} response.Response{data=dto.WechatLoginResponse}
// @Router /api/v1/auth/wechat/login [post]
func (h *UserHandler) WechatLogin(c *gin.Context) {
	var req dto.WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	result, err := h.userService.WechatLogin(req.Code)
	if err != nil {
		response.Error(c, 500, "登录失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=dto.UserInfoDTO}
// @Router /api/v1/user/info [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	userInfo, err := h.userService.GetUserInfo(userID)
	if err != nil {
		response.Error(c, 500, "获取用户信息失败: "+err.Error())
		return
	}

	response.Success(c, userInfo)
}

// UpdateUserInfo 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.UpdateUserRequest true "更新请求"
// @Success 200 {object} response.Response{data=dto.UserInfoDTO}
// @Router /api/v1/user/info [put]
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	userInfo, err := h.userService.UpdateUserInfo(userID, &req)
	if err != nil {
		response.Error(c, 500, "更新用户信息失败: "+err.Error())
		return
	}

	response.Success(c, userInfo)
}

// GenerateTestToken 生成测试 Token（仅用于开发测试）
// @Summary 生成测试 Token
// @Description 为指定用户生成测试用的 JWT Token（仅开发环境使用）
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body dto.TestTokenRequest true "测试 Token 请求"
// @Success 200 {object} response.Response{data=dto.WechatLoginResponse}
// @Router /api/v1/auth/test-token [post]
func (h *UserHandler) GenerateTestToken(c *gin.Context) {
	var req dto.TestTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	result, err := h.userService.GenerateTestToken(req.UserID)
	if err != nil {
		response.Error(c, 500, "生成 Token 失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

