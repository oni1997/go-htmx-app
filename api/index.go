package api

import (
	"go-htmx-app/handlers"
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
		templates.ExecuteTemplate(w, "wallet_status.html", nil)
	case "/cusd-balance":
		handlers.CUSDBalanceHandler(w, r)
	case "/transfer-cusd":
		handlers.TransferCUSDHandler(w, r)
	default:
		http.FileServer(http.Dir("static")).ServeHTTP(w, r)
	}
}
