package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"

	"voteweb/internal/util"
)

// ProxiedIP extracts the real client IP from X-Forwarded-For header if behind trusted proxy
func ProxiedIP(trustProxy bool, allowedCIDRs []*net.IPNet) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if trustProxy {
			// Get X-Forwarded-For header
			xff := c.GetHeader("X-Forwarded-For")
			if xff != "" {
				// Parse the first IP from X-Forwarded-For
				ips := strings.Split(xff, ",")
				if len(ips) > 0 {
					firstIP := strings.TrimSpace(ips[0])

					// Verify the request came from an allowed proxy
					remoteIP := util.NormalizeIP(c.Request.RemoteAddr)
					if util.IsIPInCIDRs(remoteIP, allowedCIDRs) {
						clientIP = firstIP
					}
				}
			}
		}

		// Normalize and store the client IP
		clientIP = util.NormalizeIP(clientIP)
		c.Set("client_ip", clientIP)

		c.Next()
	}
}


