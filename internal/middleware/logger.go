package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"magickingdom-go/internal/logger"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(startTime)

		// 获取状态码
		statusCode := c.Writer.Status()

		// 记录日志
		logger.GetLogger().WithFields(map[string]interface{}{
			"status_code": statusCode,
			"latency":     latency,
			"client_ip":   c.ClientIP(),
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
		}).Info("HTTP Request")
	}
}

