package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"gorm.io/gorm"
	"magickingdom-go/internal/config"
	"magickingdom-go/internal/dto"
	"magickingdom-go/internal/models"
	"magickingdom-go/internal/repository"
	"magickingdom-go/internal/utils"
)

// WechatSession 微信登录返回的 session 信息
type WechatSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// UserService 用户服务接口
type UserService interface {
	WechatLogin(code string) (*dto.WechatLoginResponse, error)
	GetUserInfo(userID uint) (*dto.UserInfoDTO, error)
	UpdateUserInfo(userID uint, req *dto.UpdateUserRequest) (*dto.UserInfoDTO, error)
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
	jwtUtil  *utils.JWTUtil
	cfg      *config.Config
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository, jwtUtil *utils.JWTUtil, cfg *config.Config) UserService {
	return &userService{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
		cfg:      cfg,
	}
}

// WechatLogin 微信登录
func (s *userService) WechatLogin(code string) (*dto.WechatLoginResponse, error) {
	// 调用微信接口获取 openid 和 session_key
	session, err := s.getWechatSession(code)
	if err != nil {
		return nil, fmt.Errorf("获取微信 session 失败: %w", err)
	}

	if session.ErrCode != 0 {
		return nil, fmt.Errorf("微信登录失败: %s", session.ErrMsg)
	}

	// 查找或创建用户
	user, err := s.userRepo.FindByOpenID(session.OpenID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，创建新用户
			user = &models.User{
				OpenID:     session.OpenID,
				UnionID:    session.UnionID,
				SessionKey: session.SessionKey,
				Status:     1,
			}
			if err := s.userRepo.Create(user); err != nil {
				return nil, fmt.Errorf("创建用户失败: %w", err)
			}
		} else {
			return nil, fmt.Errorf("查询用户失败: %w", err)
		}
	} else {
		// 用户已存在，更新 session_key
		user.SessionKey = session.SessionKey
		if session.UnionID != "" {
			user.UnionID = session.UnionID
		}
		if err := s.userRepo.Update(user); err != nil {
			return nil, fmt.Errorf("更新用户失败: %w", err)
		}
	}

	// 生成 JWT token
	token, err := s.jwtUtil.GenerateToken(user.ID, user.OpenID)
	if err != nil {
		return nil, fmt.Errorf("生成 token 失败: %w", err)
	}

	return &dto.WechatLoginResponse{
		Token: token,
		User:  s.modelToDTO(user),
	}, nil
}

// GetUserInfo 获取用户信息
func (s *userService) GetUserInfo(userID uint) (*dto.UserInfoDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	userDTO := s.modelToDTO(user)
	return &userDTO, nil
}

// UpdateUserInfo 更新用户信息
func (s *userService) UpdateUserInfo(userID uint, req *dto.UpdateUserRequest) (*dto.UserInfoDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}
	if req.Country != "" {
		user.Country = req.Country
	}
	if req.Province != "" {
		user.Province = req.Province
	}
	if req.City != "" {
		user.City = req.City
	}
	if req.Language != "" {
		user.Language = req.Language
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	userDTO := s.modelToDTO(user)
	return &userDTO, nil
}

// getWechatSession 调用微信接口获取 session
func (s *userService) getWechatSession(code string) (*WechatSession, error) {
	url := fmt.Sprintf("%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		s.cfg.Wechat.LoginURL,
		s.cfg.Wechat.AppID,
		s.cfg.Wechat.AppSecret,
		code,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var session WechatSession
	if err := json.Unmarshal(body, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// modelToDTO 将 model 转换为 DTO
func (s *userService) modelToDTO(user *models.User) dto.UserInfoDTO {
	return dto.UserInfoDTO{
		ID:        user.ID,
		OpenID:    user.OpenID,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
		Gender:    user.Gender,
		Country:   user.Country,
		Province:  user.Province,
		City:      user.City,
		Language:  user.Language,
		Phone:     user.Phone,
		Email:     user.Email,
		Status:    user.Status,
	}
}

