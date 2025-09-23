package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type IPLocation struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GetCurrentLocation fetches current location using IP API
func GetCurrentLocation() (string, float64, float64, error) {
	resp, err := http.Get("https://ipapi.co/json/")
	if err != nil {
		return "", 0, 0, err
	}
	defer resp.Body.Close()

	var loc IPLocation
	if err := json.NewDecoder(resp.Body).Decode(&loc); err != nil {
		return "", 0, 0, err
	}

	if loc.City == "" {
		return "", 0, 0, errors.New("could not detect location")
	}

	return fmt.Sprintf("%s,%s", loc.City, loc.Country), loc.Latitude, loc.Longitude, nil
}
