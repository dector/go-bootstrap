package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"mymyapp/internal/config"
	"mymyapp/internal/server"
	"mymyapp/meta"
)

func main() {
	moto := fmt.Sprintf("%s v%s", meta.G.Name, meta.G.Version)
	fmt.Println(moto)

	cfg := config.NewServerConfig()
	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start HTTP server: %v\n", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Printf("Running on http://%s\n", listener.Addr())
	if err := run(ctx, listener, server.New()); err != nil {
		fmt.Fprintf(os.Stderr, "HTTP server: %v\n", err)
		os.Exit(1)
	}
}
