package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"magickingdom-go/internal/logger"
	"magickingdom-go/internal/response"
)

// RecoveryMiddleware 恢复中间件，捕获 panic
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.GetLogger().WithFields(map[string]interface{}{
					"error": err,
					"path":  c.Request.URL.Path,
				}).Error("Panic recovered")

				response.Error(c, http.StatusInternalServerError, "服务器内部错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}

