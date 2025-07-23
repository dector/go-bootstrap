package main

import (
	"fmt"
	"mymyapp/internal/config"
	"mymyapp/internal/ui/components"
	"mymyapp/internal/ui/pages"
	"mymyapp/meta"
	"net/http"

	"github.com/go-chi/chi/v5"
	mw "github.com/go-chi/chi/v5/middleware"
)

func main() {
	moto := fmt.Sprintf("%s v%s", meta.G.Name, meta.G.Version)
	fmt.Println(moto)

	cfg := config.NewServerConfig()

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

	fmt.Printf("Running on http://localhost:%s\n", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}
