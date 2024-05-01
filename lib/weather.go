package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type WeatherDesc struct {
	Value string `json:"value"`
}

type WeatherDetail struct {
	TempC          string        `json:"temp_C"`
	TempF          string        `json:"temp_F"`
	Winddir16Point string        `json:"winddir16Point"`
	WindSpeedKmph  string        `json:"windspeedKmph"`
	WindSpeedMiles string        `json:"windspeedMiles"`
	WeatherDesc    []WeatherDesc `json:"weatherDesc"`
	Chanceofrain   *string       `json:"chanceofrain,omitempty"`
}

type Weather struct {
	Date   string          `json:"date"`
	Hourly []WeatherDetail `json:"hourly"`
}

type WeatherInfo struct {
	CurrentCondition []WeatherDetail `json:"current_condition"`
	Weather          []Weather       `json:"weather"`
}

func GetWeatherInfo(city string) (*WeatherInfo, error) {
	cityFormatted := strings.ReplaceAll(city, " ", "+")
	url := "https://wttr.in/" + cityFormatted + "?format=j1"

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to get weather info for city: %s", city)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var weatherInfo WeatherInfo
	err = json.Unmarshal(resBody, &weatherInfo)
	if err != nil {
		return nil, err
	}

	slog.Info("Weather info for city", "city", city, "info", fmt.Sprintf("%+v", weatherInfo))

	return &weatherInfo, nil
}
