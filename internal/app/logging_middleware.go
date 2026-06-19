package app

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func requestLoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(context *gin.Context) {
		startedAt := time.Now()

		context.Next()

		path := context.FullPath()
		if path == "" {
			path = context.Request.URL.Path
		}

		fields := []zap.Field{
			zap.String("method", context.Request.Method),
			zap.String("path", path),
			zap.Int("status", context.Writer.Status()),
			zap.Duration("duration", time.Since(startedAt)),
			zap.String("client_ip", context.ClientIP()),
		}
		if len(context.Errors) > 0 {
			fields = append(fields, zap.Strings("errors", context.Errors.Errors()))
		}

		if context.Writer.Status() >= http.StatusInternalServerError {
			logger.Error("HTTP request", fields...)
			return
		}

		logger.Info("HTTP request", fields...)
	}
}

func recoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			recovered := recover()
			if recovered == nil {
				return
			}

			_ = context.Error(fmt.Errorf("panic recovered: %v", recovered))
			logger.Error(
				"panic recovered",
				zap.String("method", context.Request.Method),
				zap.String("path", context.Request.URL.Path),
				zap.Any("panic", recovered),
				zap.ByteString("stack", debug.Stack()),
			)
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "internal_error",
					"message": "Internal server error",
					"fields":  []any{},
				},
			})
		}()

		context.Next()
	}
}
