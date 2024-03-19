package main

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

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

	// Handle static files from the embed FS (with a custom handler).
	http.Handle("GET /static/", gowebly.StaticFileServerHandler(http.FS(static)))

	// Handle index page view.
	http.HandleFunc("GET /", indexViewHandler)

	// Handle API endpoints.
	http.HandleFunc("GET /api/hello-world", showContentAPIHandler)

	// Create a new server instance with options from environment variables.
	// For more information, see https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Send log message.
	slog.Info("Starting server...", "port", port)

	return server.ListenAndServe()
}
