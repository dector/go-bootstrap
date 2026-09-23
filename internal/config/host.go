package config

import "os"

const DefaultHost = "localhost"

func parseHostOrDefault() string {
	host := os.Getenv("HOST")
	if host == "" {
		return DefaultHost
	}

	return host
}
