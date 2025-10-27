package util

import (
	"testing"
)

func TestIPHasher_HashIP(t *testing.T) {
	salt := "test-salt-for-hashing"
	hasher := NewIPHasher(salt)

	tests := []struct {
		name string
		ip1  string
		ip2  string
		want bool // true if hashes should be equal
	}{
		{
			name: "same IP produces same hash",
			ip1:  "192.168.1.1",
			ip2:  "192.168.1.1",
			want: true,
		},
		{
			name: "different IPs produce different hashes",
			ip1:  "192.168.1.1",
			ip2:  "192.168.1.2",
			want: false,
		},
		{
			name: "IPv6 addresses",
			ip1:  "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			ip2:  "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash1 := hasher.HashIP(tt.ip1)
			hash2 := hasher.HashIP(tt.ip2)

			equal := string(hash1) == string(hash2)
			if equal != tt.want {
				t.Errorf("HashIP() equal = %v, want %v", equal, tt.want)
			}

			// Hash should always be 32 bytes (SHA256)
			if len(hash1) != 32 {
				t.Errorf("HashIP() length = %v, want 32", len(hash1))
			}
		})
	}
}

func TestIPHasher_DifferentSalts(t *testing.T) {
	ip := "192.168.1.1"

	hasher1 := NewIPHasher("salt1")
	hasher2 := NewIPHasher("salt2")

	hash1 := hasher1.HashIP(ip)
	hash2 := hasher2.HashIP(ip)

	if string(hash1) == string(hash2) {
		t.Error("Different salts should produce different hashes for the same IP")
	}
}

