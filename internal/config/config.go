package config

import "strconv"

type ServerConfig struct {
	Port string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port: parsePortOrDefault(),
	}
}

func (self ServerConfig) PortInt() int {
	portInt, err := strconv.Atoi(self.Port)
	if err != nil {
		panic("Invalid port value: " + self.Port)
	}

	return portInt
}
