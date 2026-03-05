package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 声明
type Claims struct {
	UserID uint   `json:"user_id"`
	OpenID string `json:"open_id"`
	jwt.RegisteredClaims
}

// JWTUtil JWT 工具
type JWTUtil struct {
	secret []byte
	expire time.Duration
}

// NewJWTUtil 创建 JWT 工具实例
func NewJWTUtil(secret string, expire time.Duration) *JWTUtil {
	return &JWTUtil{
		secret: []byte(secret),
		expire: expire,
	}
}

// GenerateToken 生成 Token
func (j *JWTUtil) GenerateToken(userID uint, openID string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		OpenID: openID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expire)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ParseToken 解析 Token
func (j *JWTUtil) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的 token")
}

