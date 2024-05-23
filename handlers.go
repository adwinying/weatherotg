package main

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/a-h/templ"
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

	// Extract location param from query string
	location := r.URL.Query().Get("location")
	isLocationSet := location != ""
	if location == "" {
		// Get location from IP
		ip := strings.Split(r.RemoteAddr, ":")[0]
		city, err := lib.GetCityFromIp(ip)
		if err != nil || city == "Undefined" {
			city = "Tokyo"
		}
		location = city
	}

	// Extract mode param from query string
  modeQuery := r.URL.Query().Get("mode")
	mode, err := lib.ParseDisplayMode(modeQuery)
  isModeSet := modeQuery != ""
	if err != nil {
		mode = lib.Default
	}

	// Extract unit param from query string
  unitQuery := r.URL.Query().Get("unit")
	unit, err := lib.ParseTemperatureUnit(unitQuery)
  isUnitSet := unitQuery != ""
	if err != nil {
		unit = lib.Celsius
	}

	// Get weather info
	weatherInfo, err := lib.GetWeatherInfo(location)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Format weather info
	formattedWeatherInfo, err := lib.FormatWeatherInfo(weatherInfo)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Define template layout for index page.
	page := func() templ.Component {
		// Return weather info component only if request is htmx
		if htmx.IsHTMX(r) {
			return pages.IndexContent(
				isLocationSet,
				location,
        isModeSet,
				mode,
        isUnitSet,
				unit,
				formattedWeatherInfo,
			)
		}

		return templates.Layout(
			templates.MetaTags("WeatherOTG", "", ""),
			&mode,
			pages.IndexContent(
				isLocationSet,
				location,
        isModeSet,
				mode,
        isUnitSet,
				unit,
				formattedWeatherInfo,
			),
		)
	}()

	// Render index page template.
	if err := htmx.NewResponse().RenderTempl(r.Context(), w, page); err != nil {
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
		nil,
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
		nil,
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
