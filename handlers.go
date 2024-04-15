package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/angelofallars/htmx-go"

	"github.com/adwinying/weatherotg/lib"
	"github.com/adwinying/weatherotg/templates"
	"github.com/adwinying/weatherotg/templates/pages"
)

func indexViewHandler(w http.ResponseWriter, r *http.Request) {
	// Check if path is not root
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	// Get location from IP
	ip := strings.Split(r.RemoteAddr, ":")[0]
	city, err := lib.GetCityFromIp(ip)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Define template layout for index page.
	indexTemplate := templates.Layout(
		templates.MetaTags("", "", ""),
		pages.IndexContent(
			"Welcome to example!",               // define h1 text
			fmt.Sprintf("You are in %s.", city), // define p text
		),
	)

	// Render index page template.
	if err := htmx.NewResponse().RenderTempl(r.Context(), w, indexTemplate); err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Send log message.
	slog.Info("", "method", r.Method, "status", http.StatusOK, "path", r.URL.Path)
}

func aboutViewHandler(w http.ResponseWriter, r *http.Request) {
	// Define template layout for about page.
	aboutTemplate := templates.Layout(
		templates.MetaTags("About", "", ""),
		pages.AboutContent(),
	)

	// Render about page template.
	if err := htmx.NewResponse().RenderTempl(r.Context(), w, aboutTemplate); err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Send log message.
	slog.Info("", "method", r.Method, "status", http.StatusOK, "path", r.URL.Path)
}

// showContentAPIHandler handles an API endpoint to show content.
func showContentAPIHandler(w http.ResponseWriter, r *http.Request) {
	// Check, if the current request has a 'HX-Request' header.
	// For more information, see https://htmx.org/docs/#request-headers
	if !htmx.IsHTMX(r) {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Write HTML content.
	w.Write([]byte("<p>🎉 Yes, <strong>htmx</strong> is ready to use! (<code>GET /api/hello-world</code>)</p>"))

	// Send htmx response.
	htmx.NewResponse().Write(w)

	// Send log message.
	slog.Info("", "method", r.Method, "status", http.StatusOK, "path", r.URL.Path)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	errName := http.StatusText(status)
	if errName == "" {
		errName = "Error"
	}

	errDescription := "An error has occurred. Please try again later."
	if status == http.StatusNotFound {
		errDescription = "The page you are looking for does not exist."
	}

	// Define template layout for error page.
	errorTemplate := templates.Layout(
		templates.MetaTags(errName, "", ""),
		pages.ErrorContent(status, errName, errDescription),
	)

	// Set HTTP status
	w.WriteHeader(status)

	// Render error page template.
	if err := htmx.NewResponse().RenderTempl(r.Context(), w, errorTemplate); err != nil {
		slog.Error("", "method", r.Method, "status", status, "path", r.URL.Path)
		return
	}

	slog.Info("", "method", r.Method, "status", status, "path", r.URL.Path)
}
