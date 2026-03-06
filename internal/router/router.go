package router

import (
	"magickingdom-go/internal/handler"
	"magickingdom-go/internal/middleware"
	"magickingdom-go/internal/utils"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(
	userHandler *handler.UserHandler,
	seatHandler *handler.SeatHandler,
	jwtUtil *utils.JWTUtil,
) *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CORSMiddleware())

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由（无需 JWT）
		auth := v1.Group("/auth")
		{
			auth.POST("/wechat/login", userHandler.WechatLogin)
			// 测试用：生成测试 Token（生产环境应删除）
			auth.POST("/test-token", userHandler.GenerateTestToken)
		}

		// 用户相关路由（需要 JWT）
		user := v1.Group("/user")
		user.Use(middleware.AuthMiddleware(jwtUtil))
		{
			user.GET("/info", userHandler.GetUserInfo)
			user.PUT("/info", userHandler.UpdateUserInfo)
		}

		// 座位相关路由
		seats := v1.Group("/seats")
		{
			// 无需认证的接口
			seats.GET("/layout", seatHandler.GetSeatLayout)
			seats.GET("/occupancy", seatHandler.GetSeatOccupancyInfo)

			// 需要认证的接口
			seatsAuth := seats.Group("")
			seatsAuth.Use(middleware.AuthMiddleware(jwtUtil))
			{
				seatsAuth.GET("/my-seat", seatHandler.GetUserCurrentSeat)
				seatsAuth.POST("/join", seatHandler.JoinSeat)
				seatsAuth.POST("/leave", seatHandler.LeaveSeat)
			}
		}
	}

	return r
}

