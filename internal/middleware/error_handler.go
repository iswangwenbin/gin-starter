package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Error("panic recovered",
				zap.String("error", err),
				zap.String("request_id", c.GetString("request_id")),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("stack", string(debug.Stack())),
			)
			
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "internal server error",
				"data":    nil,
			})
		} else {
			logger.Error("panic recovered",
				zap.Any("error", recovered),
				zap.String("request_id", c.GetString("request_id")),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("stack", string(debug.Stack())),
			)
			
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "internal server error",
				"data":    nil,
			})
		}
		
		c.Abort()
	})
}

func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "endpoint not found",
			"data":    nil,
		})
	}
}

func MethodNotAllowedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"code":    405,
			"message": "method not allowed",
			"data":    nil,
		})
	}
}