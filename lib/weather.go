package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type WeatherDesc struct {
	Value string `json:"value"`
}

type CurrentWeatherDetail struct {
	LocalObsDateTime string        `json:"localObsDateTime"`
	TempC            string        `json:"temp_C"`
	TempF            string        `json:"temp_F"`
	Winddir16Point   string        `json:"winddir16Point"`
	WindSpeedKmph    string        `json:"windspeedKmph"`
	WindSpeedMiles   string        `json:"windspeedMiles"`
	WeatherDesc      []WeatherDesc `json:"weatherDesc"`
}

type ForecastWeatherDetail struct {
	Time           string        `json:"time"`
	TempC          string        `json:"tempC"`
	TempF          string        `json:"tempF"`
	Winddir16Point string        `json:"winddir16Point"`
	WindSpeedKmph  string        `json:"windspeedKmph"`
	WindSpeedMiles string        `json:"windspeedMiles"`
	WeatherDesc    []WeatherDesc `json:"weatherDesc"`
	ChanceOfRain   string        `json:"chanceofrain"`
}

type Weather struct {
	Date   string                  `json:"date"`
	Hourly []ForecastWeatherDetail `json:"hourly"`
}

type WeatherInfo struct {
	CurrentCondition []CurrentWeatherDetail `json:"current_condition"`
	Weather          []Weather              `json:"weather"`
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

	return &weatherInfo, nil
}

type FormattedCurrentWeatherInfo struct {
	Detail      CurrentWeatherDetail
	Description string
	WindSpeedMs float64
	Icon        string
	Timestamp   time.Time
}

type FormattedForecastWeatherInfo struct {
	Detail      ForecastWeatherDetail
	Description string
	WindSpeedMs float64
	Icon        string
	Timestamp   time.Time
}

type FormattedDailyWeatherInfo struct {
	Timestamp time.Time
	Hourly    [8]FormattedForecastWeatherInfo
}

type FormattedWeather struct {
	Current  FormattedCurrentWeatherInfo
	Forecast [4]FormattedForecastWeatherInfo
	Daily    [3]FormattedDailyWeatherInfo
}

func getWeatherIcon(description string) string {
	lowercaseDesc := strings.ToLower(description)

	switch {
	case strings.Contains(lowercaseDesc, "sunny"):
		return "mdi:weather-sunny"
	case strings.Contains(lowercaseDesc, "snow"):
		return "mdi:weather-snowy"
	case strings.Contains(lowercaseDesc, "rain"):
		return "mdi:weather-rainy"
	case strings.Contains(lowercaseDesc, "drizzle"):
		return "mdi:weather-rainy"
	case strings.Contains(lowercaseDesc, "cloudy"):
		return "mdi:weather-cloudy"
	case strings.Contains(lowercaseDesc, "clear"):
		return "mdi:weather-night"
	case strings.Contains(lowercaseDesc, "thunder"):
		return "mdi:lightning-bolt-outline"
	case strings.Contains(lowercaseDesc, "overcast"):
		return "mdi:weather-partly-cloudy"
	case strings.Contains(lowercaseDesc, "mist"):
		return "mdi:weather-fog"
	default:
		return "mdi:help-rhombus-outline"
	}
}

func getWindSpeedMs(speedKmph string) float64 {
	speed, err := strconv.ParseFloat(speedKmph, 64)
	if err != nil {
		return 0
	}

	return speed / 3.6
}

func FormatWeatherInfo(info *WeatherInfo) (*FormattedWeather, error) {
	current, err := func() (*FormattedCurrentWeatherInfo, error) {
		currentTimestamp, err := time.Parse("2006-01-02 3:04 PM", info.CurrentCondition[0].LocalObsDateTime)
		if err != nil {
			return nil, err
		}

		result := FormattedCurrentWeatherInfo{
			Detail:      info.CurrentCondition[0],
			Description: info.CurrentCondition[0].WeatherDesc[0].Value,
			WindSpeedMs: getWindSpeedMs(info.CurrentCondition[0].WindSpeedKmph),
			Icon:        getWeatherIcon(info.CurrentCondition[0].WeatherDesc[0].Value),
			Timestamp:   currentTimestamp,
		}

		return &result, nil
	}()

	if err != nil {
		return nil, err
	}

	// format daily weather info
	daily, err := func() (*[3]FormattedDailyWeatherInfo, error) {
		result := [3]FormattedDailyWeatherInfo{}

		for i := 0; i < 3; i++ {
			dailyTimestamp, err := time.Parse("2006-01-02", info.Weather[i].Date)
			if err != nil {
				return nil, err
			}

			result[i] = FormattedDailyWeatherInfo{
				Timestamp: dailyTimestamp,
			}

			for j := 0; j < 8; j++ {
				hours, err := strconv.Atoi(info.Weather[i].Hourly[j].Time)
				if err != nil {
					return nil, err
				}

				result[i].Hourly[j] = FormattedForecastWeatherInfo{
					Detail:      info.Weather[i].Hourly[j],
					Description: info.Weather[i].Hourly[j].WeatherDesc[0].Value,
					WindSpeedMs: getWindSpeedMs(info.Weather[i].Hourly[j].WindSpeedKmph),
					Icon:        getWeatherIcon(info.Weather[i].Hourly[j].WeatherDesc[0].Value),
					Timestamp:   dailyTimestamp.Add(time.Hour * time.Duration(hours/100)),
				}
			}
		}

		return &result, nil
	}()

	if err != nil {
		return nil, err
	}

	// extract the first three entries from daily starting from current timestamp
	forecast := func() [4]FormattedForecastWeatherInfo {
		combinedDaily := append(
			append(daily[0].Hourly[:], daily[1].Hourly[:]...),
			daily[2].Hourly[:]...,
		)

		firstForecastIndex := 0
		for i, hourly := range combinedDaily {
			if hourly.Timestamp.Before(current.Timestamp) {
				continue
			}

			firstForecastIndex = i
			break
		}

		result := [4]FormattedForecastWeatherInfo{}
		copy(result[:], combinedDaily[firstForecastIndex:firstForecastIndex+4])

		return result
	}()

	return &FormattedWeather{
		Current:  *current,
		Forecast: forecast,
		Daily:    *daily,
	}, nil
}
