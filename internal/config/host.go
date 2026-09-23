package config

import "os"

const DefaultHost = "localhost"

func parseHostOrDefault() string {
	host := os.Getenv("HOST")
	switch host {
	case "":
		return DefaultHost
	case "*":
		return "0.0.0.0"
	default:
		return host
	}
}
