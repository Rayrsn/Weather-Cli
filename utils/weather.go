package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Weather struct {
	Location      string  `json:"location"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Temp          float64 `json:"temperature"`
	FeelsLike     float64 `json:"feels_like"`
	MinTemp       float64 `json:"min_temp"`
	MaxTemp       float64 `json:"max_temp"`
	Humidity      float64 `json:"humidity"`
	WindSpeed     float64 `json:"wind_speed"`
	WindDirection float64 `json:"wind_direction"`
	Pressure      float64 `json:"pressure"`
	UVIndex       float64 `json:"uv_index"`
	Description   string  `json:"description"`
}

var GeocodingUrl = "https://geocoding-api.open-meteo.com/v1/search"
var ForecastUrl = "https://api.open-meteo.com/v1/forecast"
var AirQualityUrl = "https://air-quality-api.open-meteo.com/v1/air-quality"

// ----- Weather helpers -----
func translateWeatherCode(code int) string {
	switch code {
	case 0:
		return "Clear Sky"
	case 1:
		return "Mainly Clear"
	case 2:
		return "Partly Cloudy"
	case 3:
		return "Overcast"
	case 45:
		return "Fog"
	case 48:
		return "Depositing Rime Fog"
	case 51:
		return "Light Drizzle"
	case 53:
		return "Moderate Drizzle"
	case 55:
		return "Dense Drizzle"
	case 61:
		return "Slight Rain"
	case 63:
		return "Moderate Rain"
	case 65:
		return "Heavy Rain"
	case 71:
		return "Snow Fall"
	case 95:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}

func toFloat(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	default:
		return 0
	}
}

func toFloatSlice(v interface{}) []float64 {
	out := []float64{}
	if arr, ok := v.([]interface{}); ok {
		for _, x := range arr {
			out = append(out, toFloat(x))
		}
	}
	return out
}

// ----- Core Functions -----

func GetWeatherByCoords(lat, lon float64, locationName string) (Weather, error) {
	url := fmt.Sprintf("%s?latitude=%.4f&longitude=%.4f&timezone=auto&current_weather=true&hourly=temperature_2m,apparent_temperature,relativehumidity_2m,pressure_msl", ForecastUrl, lat, lon)

	resp, err := http.Get(url)
	if err != nil {
		return Weather{}, err
	}
	defer resp.Body.Close()

	var forecastData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&forecastData); err != nil {
		return Weather{}, err
	}

	current := forecastData["current_weather"].(map[string]interface{})
	temp := toFloat(current["temperature"])
	windSpeed := toFloat(current["windspeed"])
	windDir := toFloat(current["winddirection"])
	weatherCode := int(toFloat(current["weathercode"]))

	// hourly data
	hourly := forecastData["hourly"].(map[string]interface{})
	temps := toFloatSlice(hourly["temperature_2m"])
	apparent := toFloatSlice(hourly["apparent_temperature"])
	humidityArr := toFloatSlice(hourly["relativehumidity_2m"])
	pressureArr := toFloatSlice(hourly["pressure_msl"])

	hour := time.Now().Hour()
	feelsLike := apparent[hour]
	humidity := humidityArr[hour]
	pressure := pressureArr[hour]

	return Weather{
		Location:      locationName,
		Lat:           lat,
		Lon:           lon,
		Temp:          temp,
		FeelsLike:     feelsLike,
		MinTemp:       temps[hour],
		MaxTemp:       temps[hour],
		Humidity:      humidity,
		WindSpeed:     windSpeed,
		WindDirection: windDir,
		Pressure:      pressure,
		Description:   translateWeatherCode(weatherCode),
	}, nil
}

func GetWeatherByCity(city string) (Weather, error) {
	geoUrl := fmt.Sprintf("%s?name=%s&count=1", GeocodingUrl, url.QueryEscape(city))
	resp, err := http.Get(geoUrl)
	if err != nil {
		return Weather{}, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return Weather{}, err
	}

	results := data["results"].([]interface{})
	first := results[0].(map[string]interface{})
	lat := toFloat(first["latitude"])
	lon := toFloat(first["longitude"])
	name := fmt.Sprintf("%s,%s", first["name"], first["country"])

	return GetWeatherByCoords(lat, lon, name)
}
