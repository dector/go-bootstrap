package main

import (
	"fmt"
	"mymyapp/internal/config"
	"mymyapp/meta"
	"net/http"
	"strings"

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
		sb := strings.Builder{}
		sb.WriteString("<body>")
		sb.WriteString(moto)
		sb.WriteString("</body>")
		w.Write([]byte(sb.String()))
	})

	fmt.Printf("Running on http://localhost:%s\n", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}
