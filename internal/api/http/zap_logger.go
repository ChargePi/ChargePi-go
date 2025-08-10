package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// loggingMiddleware returns a gin.HandlerFunc for logging HTTP requests using zap
func loggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		zap.L().Info("HTTP Request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.String("ip", param.ClientIP),
			zap.String("user-agent", param.Request.UserAgent()),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("error", param.ErrorMessage),
		)
		return ""
	})
}
