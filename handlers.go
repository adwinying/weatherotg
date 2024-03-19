package main

import (
	"log/slog"
	"net/http"

	"github.com/angelofallars/htmx-go"

	"github.com/adwinying/weatherotg/templates"
	"github.com/adwinying/weatherotg/templates/pages"
)

// indexViewHandler handles a view for the index page.
func indexViewHandler(w http.ResponseWriter, r *http.Request) {
	// Check if path is not root
	if r.URL.Path != "/" {
		notFoundHandler(w, r)
		return
	}

	// Define template layout for index page.
	indexTemplate := templates.Layout(
		templates.MetaTags("", "", ""),
		pages.IndexContent(
			"Welcome to example!",                // define h1 text
			"You're here because it worked out.", // define p text
		),
	)

	// Render index page template.
	if err := htmx.NewResponse().RenderTempl(r.Context(), w, indexTemplate); err != nil {
		// If not, return HTTP 400 error.
		slog.Error("render template", "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send log message.
	slog.Info("render page", "method", r.Method, "status", http.StatusOK, "path", r.URL.Path)
}

// showContentAPIHandler handles an API endpoint to show content.
func showContentAPIHandler(w http.ResponseWriter, r *http.Request) {
	// Check, if the current request has a 'HX-Request' header.
	// For more information, see https://htmx.org/docs/#request-headers
	if !htmx.IsHTMX(r) {
		// If not, return HTTP 400 error.
		slog.Error("request API", "method", r.Method, "status", http.StatusBadRequest, "path", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Write HTML content.
	w.Write([]byte("<p>🎉 Yes, <strong>htmx</strong> is ready to use! (<code>GET /api/hello-world</code>)</p>"))

	// Send htmx response.
	htmx.NewResponse().Write(w)

	// Send log message.
	slog.Info("request API", "method", r.Method, "status", http.StatusOK, "path", r.URL.Path)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Define template layout for 404 page.
	notFoundTemplate := templates.Layout(
		templates.MetaTags("404 Not Found", "", ""),
		pages.NotFoundContent(),
	)

  // Set HTTP 404 status
  w.WriteHeader(http.StatusNotFound)

	// Render 404 page template.
	if err := htmx.NewResponse().RenderTempl(r.Context(), w, notFoundTemplate); err != nil {
		// If not, return HTTP 400 error.
		slog.Error("render template", "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

  slog.Info("render page", "method", r.Method, "status", http.StatusNotFound, "path", r.URL.Path)
}
