package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRunWaitsForActiveRequestOnShutdown(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	entered := make(chan struct{})
	release := make(chan struct{})
	defer func() {
		select {
		case <-release:
		default:
			close(release)
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	finished := make(chan error, 1)
	go func() {
		finished <- run(ctx, listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			close(entered)
			<-release
			w.WriteHeader(http.StatusNoContent)
		}))
	}()

	client := &http.Client{Timeout: 3 * time.Second}
	response := make(chan error, 1)
	go func() {
		resp, err := client.Get("http://" + listener.Addr().String())
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode != http.StatusNoContent {
				err = errors.New("unexpected response status")
			}
		}
		response <- err
	}()

	select {
	case <-entered:
	case <-time.After(3 * time.Second):
		t.Fatal("request never reached server")
	}
	cancel()
	select {
	case err := <-finished:
		t.Fatalf("server stopped before request completed: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
	close(release)
	select {
	case err := <-response:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("request did not finish")
	}
	select {
	case err := <-finished:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("server did not shut down")
	}
}

func TestRunReportsServeError(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}

	err = run(context.Background(), listener, http.NewServeMux())
	if err == nil || !strings.Contains(err.Error(), "serve:") {
		t.Fatalf("expected serve error, got %v", err)
	}
}
