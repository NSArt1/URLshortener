package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/you/url-shortener/internal/handler"
    "github.com/you/url-shortener/internal/storage"
)

func main() {
    dsn := os.Getenv("DATABASE_URL") // e.g. "host=... user=... password=... dbname=... sslmode=disable"
    db, err := storage.NewPostgres(dsn)
    if err != nil {
        log.Fatalf("failed to connect db: %v", err)
    }

    r := gin.Default()
    h := handler.New(db)

    // REST API
    r.POST("/shorten", h.Shorten)
    r.GET("/stats/:code", h.Stats)

    // Redirect endpoint
    r.GET("/:code", h.Redirect)

    // Prometheus metrics endpoint
    r.GET("/metrics", gin.WrapH(promhttp.Handler()))

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    r.Run(":" + port)
}