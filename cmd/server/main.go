package main

import (
	"fmt"
	"mymyapp/internal/config"
	"mymyapp/internal/ui/components"
	"mymyapp/internal/ui/pages"
	"mymyapp/internal/utils"
	"mymyapp/meta"
	"net/http"

	"github.com/go-chi/chi/v5"
	mw "github.com/go-chi/chi/v5/middleware"
)

func main() {
	moto := fmt.Sprintf("%s v%s", meta.G.Name, meta.G.Version)
	fmt.Println(moto)

	cfg := config.NewServerConfig()

	r := newRouter()

	addr := "127.0.0.1:" + cfg.Port
	fmt.Printf("Running on http://%s\n", addr)
	http.ListenAndServe(addr, r)
}

func newRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(mw.Logger)
	r.Use(mw.Recoverer)
	r.Use(mw.RealIP)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
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
		component.Render(r.Context(), w)
	})

	r.Get("/health", healthHandler)

	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Vary", "Accept")

	if utils.AcceptsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
}
