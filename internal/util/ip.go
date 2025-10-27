package util

import (
	"net"
	"strings"
)

// NormalizeIP normalizes an IP address by removing port if present
func NormalizeIP(ip string) string {
	// Remove port if present
	if strings.Contains(ip, ":") {
		host, _, err := net.SplitHostPort(ip)
		if err == nil {
			return host
		}
	}
	return ip
}

// IsIPInCIDRs checks if an IP is in any of the given CIDR ranges
func IsIPInCIDRs(ip string, cidrs []*net.IPNet) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, cidr := range cidrs {
		if cidr.Contains(parsedIP) {
			return true
		}
	}
	return false
}


