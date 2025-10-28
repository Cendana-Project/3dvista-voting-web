package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	IPHashSalt        string
	TrustProxy        bool
	AllowedProxyCIDRs []*net.IPNet
	AppBaseURL        string
	Port              string
	GinMode           string
	AdminCode         string
	AutoVoteOnView    bool
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/voteweb?sslmode=disable"),
		IPHashSalt:     getEnv("IP_HASH_SALT", ""),
		TrustProxy:     getEnvBool("TRUST_PROXY", false),
		AppBaseURL:     getEnv("APP_BASE_URL", "http://localhost:8080"),
		Port:           getEnv("PORT", "8080"),
		GinMode:        getEnv("GIN_MODE", "debug"),
		AdminCode:      getEnv("ADMIN_CODE", ""),
		AutoVoteOnView: getEnvBool("AUTO_VOTE_ON_VIEW", false),
	}

	// Validate required fields
	if cfg.IPHashSalt == "" {
		return nil, fmt.Errorf("IP_HASH_SALT is required")
	}

	// Parse allowed proxy CIDRs
	if cfg.TrustProxy {
		cidrsStr := getEnv("ALLOWED_PROXY_CIDRS", "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16")
		cidrs := strings.Split(cidrsStr, ",")
		for _, cidr := range cidrs {
			cidr = strings.TrimSpace(cidr)
			if cidr == "" {
				continue
			}
			_, ipNet, err := net.ParseCIDR(cidr)
			if err != nil {
				return nil, fmt.Errorf("invalid CIDR %s: %w", cidr, err)
			}
			cfg.AllowedProxyCIDRs = append(cfg.AllowedProxyCIDRs, ipNet)
		}
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}
