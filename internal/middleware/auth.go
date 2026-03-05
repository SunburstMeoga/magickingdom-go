package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"magickingdom-go/internal/response"
	"magickingdom-go/internal/utils"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware(jwtUtil *utils.JWTUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 中获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		// 检查 token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "认证令牌格式错误")
			c.Abort()
			return
		}

		// 解析 token
		claims, err := jwtUtil.ParseToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "认证令牌无效或已过期")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("open_id", claims.OpenID)

		c.Next()
	}
}

// GetUserID 从上下文中获取用户 ID
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetOpenID 从上下文中获取 OpenID
func GetOpenID(c *gin.Context) (string, bool) {
	openID, exists := c.Get("open_id")
	if !exists {
		return "", false
	}
	return openID.(string), true
}

