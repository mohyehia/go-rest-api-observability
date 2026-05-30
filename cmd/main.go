package main

import (
	"log"
	"net/http"
	"time"

	"github.com/mohyehia/rest-api-observability/internal/core"
	"github.com/mohyehia/rest-api-observability/internal/health"
	"github.com/mohyehia/rest-api-observability/internal/post"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	// Create prometheus registry
	registry := prometheus.NewRegistry()

	mux := http.NewServeMux()

	postsClient := post.NewPostsClient(&http.Client{
		Timeout: 5 * time.Second,
	})

	appMetrics := core.NewApplicationMetrics(registry)

	// Register prometheus metrics endpoint
	mux.Handle("GET /metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	// Register a health check endpoint
	mux.HandleFunc("GET /health", health.Handler)

	// Register posts handler
	post.RegisterHandlers(mux, postsClient, appMetrics)

	server := &http.Server{
		Addr:         ":9091",
		Handler:      mux,
		TLSConfig:    nil,
		ReadTimeout:  5 * time.Second,   // Max time to read the request body
		WriteTimeout: 10 * time.Second,  // Max time to write the response
		IdleTimeout:  120 * time.Second, // Max time to keep a Keep-Alive connection open
	}
	log.Printf("Server is starting on port %s...", server.Addr)
	err := server.ListenAndServe()

	if err != nil {
		log.Printf("Error starting server: %v\n", err)
	}
}
