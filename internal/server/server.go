package server

import (
	"net/http"
	"strings"

	"mymyapp/internal/ui/components"
	"mymyapp/internal/ui/pages"
	"mymyapp/internal/utils"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// New builds the HTTP handler and its routes.
func New() http.Handler {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(realIP)

	e.GET("/", func(c *echo.Context) error {
		r := c.Request()
		data := components.RequestData{
			Method:     r.Method,
			Path:       r.URL.Path,
			Query:      r.URL.Query(),
			Headers:    r.Header,
			Cookies:    r.Cookies(),
			RemoteAddr: r.RemoteAddr,
			UserAgent:  r.UserAgent(),
		}

		component := pages.RequestDebugPage(data)
		return component.Render(r.Context(), c.Response())
	})

	e.GET("/health", healthHandler)

	return e
}

// realIP mirrors chi's RealIP middleware, preferring the first forwarded IP.
func realIP(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		r := c.Request()
		if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			r.RemoteAddr = strings.TrimSpace(strings.Split(forwardedFor, ",")[0])
		} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			r.RemoteAddr = realIP
		}
		c.SetRequest(r)
		return next(c)
	}
}

func healthHandler(c *echo.Context) error {
	r := c.Request()
	c.Response().Header().Set("Vary", "Accept")

	if utils.AcceptsJSON(r) {
		c.Response().Header().Set("Content-Type", "application/json")
		c.Response().WriteHeader(http.StatusOK)
		_, err := c.Response().Write([]byte(`{"status":"ok"}`))
		return err
	}

	c.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Response().WriteHeader(http.StatusOK)
	_, err := c.Response().Write([]byte("OK"))
	return err
}
