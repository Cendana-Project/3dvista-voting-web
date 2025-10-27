package util

import (
	"net"
	"testing"
)

func TestNormalizeIP(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "IP without port",
			input: "192.168.1.1",
			want:  "192.168.1.1",
		},
		{
			name:  "IP with port",
			input: "192.168.1.1:8080",
			want:  "192.168.1.1",
		},
		{
			name:  "IPv6 without port",
			input: "2001:0db8:85a3::8a2e:0370:7334",
			want:  "2001:0db8:85a3::8a2e:0370:7334",
		},
		{
			name:  "IPv6 with port",
			input: "[2001:0db8:85a3::8a2e:0370:7334]:8080",
			want:  "2001:0db8:85a3::8a2e:0370:7334",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeIP(tt.input)
			if got != tt.want {
				t.Errorf("NormalizeIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIPInCIDRs(t *testing.T) {
	// Parse test CIDRs
	_, cidr1, _ := net.ParseCIDR("10.0.0.0/8")
	_, cidr2, _ := net.ParseCIDR("172.16.0.0/12")
	_, cidr3, _ := net.ParseCIDR("192.168.0.0/16")
	cidrs := []*net.IPNet{cidr1, cidr2, cidr3}

	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{
			name: "IP in first CIDR",
			ip:   "10.0.1.1",
			want: true,
		},
		{
			name: "IP in second CIDR",
			ip:   "172.16.0.1",
			want: true,
		},
		{
			name: "IP in third CIDR",
			ip:   "192.168.1.1",
			want: true,
		},
		{
			name: "IP not in any CIDR",
			ip:   "8.8.8.8",
			want: false,
		},
		{
			name: "invalid IP",
			ip:   "invalid",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsIPInCIDRs(tt.ip, cidrs)
			if got != tt.want {
				t.Errorf("IsIPInCIDRs() = %v, want %v", got, tt.want)
			}
		})
	}
}

