package handler

import (
    "context"
    "go-htmx-app/handlers"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
    "html/template"
)

var (
    templates = template.Must(template.ParseGlob("templates/*.html"))
)

func main() {
    server := &http.Server{
        Addr:    ":8080",
        Handler: setupRoutes(),
    }

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)

    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Could not listen on :8080: %v\n", err)
        }
    }()
    log.Println("Server is ready to handle requests at :8080")

    <-stop

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    log.Println("Shutting down the server...")
    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Could not gracefully shut down the server: %v\n", err)
    }
    log.Println("Server stopped")
}

func setupRoutes() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/", indexHandler)
    mux.HandleFunc("/wallet-status", WalletStatusHandler)
    mux.HandleFunc("/cusd-balance", handlers.CUSDBalanceHandler)
    mux.HandleFunc("/transfer-cusd", handlers.TransferCUSDHandler)
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    templates.ExecuteTemplate(w, "index.html", nil)
}

func WalletStatusHandler(w http.ResponseWriter, r *http.Request) {
    templates.ExecuteTemplate(w, "wallet_status.html", nil)
}
