package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security-related HTTP headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'; frame-ancestors 'none'; base-uri 'self'")

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Referrer policy
		c.Header("Referrer-Policy", "no-referrer")

		// Frame options
		c.Header("X-Frame-Options", "DENY")

		// Strict Transport Security (enable only with HTTPS)
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		c.Next()
	}
}


