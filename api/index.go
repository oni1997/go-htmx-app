package api

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var templates = template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		templates.ExecuteTemplate(w, "index.html", nil)
	case "/wallet-status":
		// Handle wallet status
	case "/cusd-balance":
		// Handle cUSD balance
	case "/transfer-cusd":
		// Handle cUSD transfer
	default:
		http.FileServer(http.Dir("static")).ServeHTTP(w, r)
	}
}