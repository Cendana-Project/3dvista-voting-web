package util

import (
	"crypto/hmac"
	"crypto/sha256"
)

// IPHasher provides IP hashing functionality using HMAC-SHA256
type IPHasher struct {
	salt []byte
}

// NewIPHasher creates a new IPHasher with the given salt
func NewIPHasher(salt string) *IPHasher {
	return &IPHasher{
		salt: []byte(salt),
	}
}

// HashIP hashes an IP address using HMAC-SHA256
func (h *IPHasher) HashIP(ip string) []byte {
	mac := hmac.New(sha256.New, h.salt)
	mac.Write([]byte(ip))
	return mac.Sum(nil)
}


