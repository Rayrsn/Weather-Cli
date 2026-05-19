package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCityInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":[{"name":"Test City","country":"Test Country","latitude":1.23,"longitude":4.56,"timezone":"UTC","population":1000}]}`))
	}))
	defer server.Close()

	client := NewClient()
	client.GeocodingUrl = server.URL

	resp, err := client.GetCityInfo("Test")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(resp.Results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(resp.Results))
	}

	if resp.Results[0].Name != "Test City" {
		t.Errorf("Expected Test City, got %s", resp.Results[0].Name)
	}
}

func TestGetForecast(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current_weather":{"temperature":20.5,"windspeed":10.0,"winddirection":180.0,"weathercode":0},"hourly":{"relativehumidity_2m":[50.0],"apparent_temperature":[20.0],"surface_pressure":[1013.0],"pressure_msl":[1013.0]}}`))
	}))
	defer server.Close()

	client := NewClient()
	client.ForecastUrl = server.URL

	resp, err := client.GetForecast(1.23, 4.56, false, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.CurrentWeather.Temperature != 20.5 {
		t.Errorf("Expected 20.5, got %v", resp.CurrentWeather.Temperature)
	}
}
