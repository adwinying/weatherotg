package main

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	gowebly "github.com/gowebly/helpers"
)

//go:embed all:static
var static embed.FS

// runServer runs a new HTTP server with the loaded environment variables.
func runServer() error {
	// Validate environment variables.
	port, err := strconv.Atoi(gowebly.Getenv("BACKEND_PORT", "6969"))
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /up", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("👋 Up and running!"))
	})

	// Handle static files from the embed FS (with a custom handler).
	mux.Handle("GET /static/", gowebly.StaticFileServerHandler(http.FS(static)))

	// Handle page views
	mux.HandleFunc("GET /", indexViewHandler)
	mux.HandleFunc("GET /about", aboutViewHandler)

	// Handle API endpoints.
	mux.HandleFunc("GET /api/hello-world", showContentAPIHandler)

	// Send log message.
	slog.Info("Starting server...", "port", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
