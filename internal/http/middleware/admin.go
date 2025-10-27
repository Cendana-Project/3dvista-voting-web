package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const AdminHeaderKey = "X-ADMIN-CODE"

// AdminAuth validates the X-ADMIN-CODE header for admin access
func AdminAuth(validCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		providedCode := c.GetHeader(AdminHeaderKey)

		if providedCode == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "X-ADMIN-CODE header is required",
			})
			c.Abort()
			return
		}

		// Trim whitespace for both codes
		providedCode = strings.TrimSpace(providedCode)
		validCode = strings.TrimSpace(validCode)

		if providedCode != validCode {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid admin code",
			})
			c.Abort()
			return
		}

		// Set admin flag in context for use in handlers
		c.Set("is_admin", true)

		c.Next()
	}
}
