package middleware

import (
	"github.com/gin-gonic/gin"
)

func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 安全头设置
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		// 移除服务器标识
		c.Header("Server", "")
		
		c.Next()
	}
}

func HidePoweredBy() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Powered-By", "")
		c.Next()
	}
}