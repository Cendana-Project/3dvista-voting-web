package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	csrfCookieName = "csrf_token"
	csrfHeaderName = "X-CSRF-Token"
)

// CSRF implements double-submit cookie pattern for CSRF protection
func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		// For GET requests, generate and set CSRF token
		if c.Request.Method == http.MethodGet {
			// Check if cookie already exists
			_, err := c.Cookie(csrfCookieName)
			if err != nil {
				// Generate new token
				token := generateCSRFToken()
				c.SetCookie(
					csrfCookieName,
					token,
					3600,           // 1 hour
					"/",            // path
					"",             // domain (empty means current domain)
					false,          // secure (set to true in production with HTTPS)
					true,           // httpOnly
				)
				c.Set("csrf_token", token)
			}
			c.Next()
			return
		}

		// For POST, PUT, DELETE, etc., verify CSRF token
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead && c.Request.Method != http.MethodOptions {
			cookieToken, err := c.Cookie(csrfCookieName)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "CSRF token missing",
				})
				return
			}

			headerToken := c.GetHeader(csrfHeaderName)
			if headerToken == "" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "CSRF token missing in header",
				})
				return
			}

			if cookieToken != headerToken {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "CSRF token mismatch",
				})
				return
			}
		}

		c.Next()
	}
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// GetCSRFToken returns the CSRF token from the context or cookie
func GetCSRFToken(c *gin.Context) string {
	// First try to get from context (freshly generated)
	if token, exists := c.Get("csrf_token"); exists {
		if tokenStr, ok := token.(string); ok {
			return tokenStr
		}
	}

	// Fall back to cookie
	token, err := c.Cookie(csrfCookieName)
	if err == nil {
		return token
	}

	return ""
}




