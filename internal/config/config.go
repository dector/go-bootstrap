package config

import "strconv"

type ServerConfig struct {
	Host string
	Port string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Host: parseHostOrDefault(),
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
