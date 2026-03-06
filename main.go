package main

import (
	"fmt"
	"log"
	"os"

	"magickingdom-go/internal/config"
	"magickingdom-go/internal/database"
	"magickingdom-go/internal/handler"
	"magickingdom-go/internal/logger"
	"magickingdom-go/internal/repository"
	"magickingdom-go/internal/router"
	"magickingdom-go/internal/service"
	"magickingdom-go/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	if err := logger.Init(&cfg.Log); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 初始化数据库
	db, err := database.InitDB(cfg)
	if err != nil {
		logger.GetLogger().Fatalf("初始化数据库失败: %v", err)
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化 JWT 工具
	jwtUtil := utils.NewJWTUtil(cfg.JWT.Secret, cfg.GetJWTExpireDuration())

	// 依赖注入 - 初始化各层
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtUtil, cfg)
	userHandler := handler.NewUserHandler(userService)

	// 初始化座位服务
	seatService := service.NewSeatService()
	seatOccupancyRepo := repository.NewSeatOccupancyRepository(db)
	seatOccupancyService := service.NewSeatOccupancyService(seatOccupancyRepo)
	seatHandler := handler.NewSeatHandler(seatService, seatOccupancyService)

	// 设置路由
	r := router.SetupRouter(userHandler, seatHandler, jwtUtil)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.GetLogger().Infof("服务器启动在 %s", addr)

	if err := r.Run(addr); err != nil {
		logger.GetLogger().Fatalf("服务器启动失败: %v", err)
	}
}
