package main

import (
	"fmt"
	"net"
	"net/http"

	"mymyapp/internal/config"
	"mymyapp/internal/server"
	"mymyapp/meta"
)

func main() {
	moto := fmt.Sprintf("%s v%s", meta.G.Name, meta.G.Version)
	fmt.Println(moto)

	cfg := config.NewServerConfig()
	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	fmt.Printf("Running on http://%s\n", addr)
	http.ListenAndServe(addr, server.New())
}
