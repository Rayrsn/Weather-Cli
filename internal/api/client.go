package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	HttpClient    *http.Client
	GeocodingUrl  string
	ForecastUrl   string
	AirQualityUrl string
}

func NewClient() *Client {
	return &Client{
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		GeocodingUrl:  "https://geocoding-api.open-meteo.com/v1/search",
		ForecastUrl:   "https://api.open-meteo.com/v1/forecast",
		AirQualityUrl: "https://air-quality-api.open-meteo.com/v1/air-quality",
	}
}

func (c *Client) GetCityInfo(cityName string) (*GeocodingResponse, error) {
	cityinfoUrl := fmt.Sprintf("%s?name=%s&count=10", c.GeocodingUrl, url.QueryEscape(cityName))
	resp, err := c.HttpClient.Get(cityinfoUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to reach geocoding API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding API returned status: %s", resp.Status)
	}

	var data GeocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode geocoding response: %w", err)
	}

	return &data, nil
}

func (c *Client) GetForecast(lat, lon float64, imperial bool, showForecast bool) (*ForecastResponse, error) {
	forecastUrl := fmt.Sprintf("%s?timezone=auto&latitude=%.4f&longitude=%.4f&current_weather=true&hourly=temperature_2m,weathercode,relativehumidity_2m,apparent_temperature,surface_pressure,pressure_msl", c.ForecastUrl, lat, lon)
	if showForecast {
		forecastUrl += "&daily=temperature_2m_max,temperature_2m_min,weathercode,sunrise,sunset,uv_index_max,precipitation_sum,windspeed_10m_max,winddirection_10m_dominant"
	}
	if imperial {
		forecastUrl += "&temperature_unit=fahrenheit&windspeed_unit=mph"
	}

	resp, err := c.HttpClient.Get(forecastUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to reach forecast API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forecast API returned status: %s", resp.Status)
	}

	var data ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	return &data, nil
}

func (c *Client) GetAirQuality(lat, lon float64) (*AirQualityResponse, error) {
	airqualityUrl := fmt.Sprintf("%s?timezone=auto&latitude=%.4f&longitude=%.4f&hourly=uv_index", c.AirQualityUrl, lat, lon)
	resp, err := c.HttpClient.Get(airqualityUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to reach air quality API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("air quality API returned status: %s", resp.Status)
	}

	var data AirQualityResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode air quality response: %w", err)
	}

	return &data, nil
}

func (c *Client) GetAutoLocation() (string, error) {
	resp, err := c.HttpClient.Get("http://ip-api.com/json/")
	if err != nil {
		return "", fmt.Errorf("failed to reach geolocation API: %w", err)
	}
	defer resp.Body.Close()

	var data struct {
		City string `json:"city"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode geolocation response: %w", err)
	}

	if data.City == "" {
		return "", fmt.Errorf("could not determine city from IP")
	}

	return data.City, nil
}
