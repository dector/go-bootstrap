package config

import (
	"os"
	"strconv"
)

const DefaultPort = "8080"

func parsePortOrDefault() string {
	port := os.Getenv("PORT")
	if isValidPort(port) {
		return port
	}

	return DefaultPort
}

func isValidPort(port string) bool {
	_, err := strconv.Atoi(port)
	return err == nil
}
