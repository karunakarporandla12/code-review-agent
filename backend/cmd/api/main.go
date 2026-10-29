package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "code-review-agent/backend/internal/app"
    "code-review-agent/backend/internal/config"
    "code-review-agent/backend/internal/database"
)

func main() {
    cfg := config.Load()
    db, err := database.Open(context.Background(), cfg.DatabaseURL)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    application := app.New(cfg, db)
    server := &http.Server{
        Addr:              ":" + cfg.Port,
        Handler:           application.Routes(),
        ReadHeaderTimeout: 10 * time.Second,
        ReadTimeout:       30 * time.Second,
        WriteTimeout:      60 * time.Second,
        IdleTimeout:       120 * time.Second,
    }

    go func() {
        log.Printf("API running at http://localhost:%s", cfg.Port)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
    <-stop

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    _ = server.Shutdown(ctx)
}
