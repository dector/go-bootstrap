package main

import (
	"fmt"
	"mymyapp/internal/config"
	"mymyapp/meta"
	"net/http"
	"strings"
)

func main() {
	moto := fmt.Sprintf("%s v%s", meta.G.Name, meta.G.Version)
	fmt.Println(moto)

	cfg := config.NewServerConfig()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sb := strings.Builder{}
		sb.WriteString("<body>")
		sb.WriteString(moto)
		sb.WriteString("</body>")
		w.Write([]byte(sb.String()))
	})

	fmt.Printf("Running on http://localhost:%s\n", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, http.DefaultServeMux)
}
