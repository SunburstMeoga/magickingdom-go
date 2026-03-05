package router

import (
	"github.com/gin-gonic/gin"
	"magickingdom-go/internal/handler"
	"magickingdom-go/internal/middleware"
	"magickingdom-go/internal/utils"
)

// SetupRouter 设置路由
func SetupRouter(
	userHandler *handler.UserHandler,
	jwtUtil *utils.JWTUtil,
) *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CORSMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由（无需 JWT）
		auth := v1.Group("/auth")
		{
			auth.POST("/wechat/login", userHandler.WechatLogin)
		}

		// 用户相关路由（需要 JWT）
		user := v1.Group("/user")
		user.Use(middleware.AuthMiddleware(jwtUtil))
		{
			user.GET("/info", userHandler.GetUserInfo)
			user.PUT("/info", userHandler.UpdateUserInfo)
		}
	}

	return r
}

