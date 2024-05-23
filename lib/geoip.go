package lib

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func GetCityFromIp(ip string) (string, error) {
	url := fmt.Sprintf("https://ipapi.co/%s/city/", ip)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("Failed to get city from IP: %s", ip)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	slog.Info("City from IP", "ip", ip, "city", string(resBody))

	return string(resBody), nil
}
