package config

import (
	"net"
	"testing"
)

func TestServerConfigHost(t *testing.T) {
	for _, tt := range []struct {
		name string
		host string
		want string
	}{
		{"empty", "", "localhost"},
		{"explicit IPv4", "0.0.0.0", "0.0.0.0"},
		{"explicit IPv6", "::1", "::1"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HOST", tt.host)
			t.Setenv("PORT", "9090")
			cfg := NewServerConfig()
			if cfg.Host != tt.want {
				t.Errorf("Host = %q, want %q", cfg.Host, tt.want)
			}
			if addr := net.JoinHostPort(cfg.Host, cfg.Port); addr != net.JoinHostPort(tt.want, "9090") {
				t.Errorf("address = %q, want %q", addr, net.JoinHostPort(tt.want, "9090"))
			}
		})
	}
}
