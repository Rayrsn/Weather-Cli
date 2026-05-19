/*
Copyright © 2022 Rayr https://rayr.ml/LinkInBio/
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

var GeocodingUrl = "https://geocoding-api.open-meteo.com/v1/search"
var ForecastUrl = "https://api.open-meteo.com/v1/forecast"
var AirQualityUrl = "https://air-quality-api.open-meteo.com/v1/air-quality"

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets the weather for a city",
	Long:  `Gets the weather info for a city. (Can be used with --raw to get a json response)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("please enter a city name")
		}
		var CityNameFormatted = strings.Replace(args[0], " ", "%20", -1)
		var CityName = args[0]

		if cmd.Flag("raw").Value.String() == "false" {
			fmt.Printf("Searching for city %s...\n\n", strings.ToUpper(CityName[:1])+CityName[1:])
		}
		cityinfoUrl := GeocodingUrl + "?name=" + CityNameFormatted + "&count=1"

		resp, err := http.Get(cityinfoUrl)
		if err != nil {
			return fmt.Errorf("failed to reach geocoding API: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("geocoding API returned status: %s", resp.Status)
		}

		var cityinfoData GeocodingResponse
		if err := json.NewDecoder(resp.Body).Decode(&cityinfoData); err != nil {
			return fmt.Errorf("failed to decode geocoding response: %w", err)
		}

		if len(cityinfoData.Results) == 0 {
			return fmt.Errorf("city not found: %s", CityName)
		}

		result := cityinfoData.Results[0]
		var FetchedCityName = result.Name
		var FetchedCountryName = result.Country
		var FetchedLatitude = result.Latitude
		var FetchedLongitude = result.Longitude
		var FetchedTimezone = result.Timezone
		var FetchedPopulationInt = int64(result.Population)

		forecastUrl := ForecastUrl + "?timezone=auto" + "&latitude=" + fmt.Sprintf("%.4f", FetchedLatitude) + "&longitude=" + fmt.Sprintf("%.4f", FetchedLongitude) + "&current_weather=true" + "&hourly=relativehumidity_2m,apparent_temperature,surface_pressure,pressure_msl"
		resp, err = http.Get(forecastUrl)
		if err != nil {
			return fmt.Errorf("failed to reach forecast API: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("forecast API returned status: %s", resp.Status)
		}

		var forecastData ForecastResponse
		if err := json.NewDecoder(resp.Body).Decode(&forecastData); err != nil {
			return fmt.Errorf("failed to decode forecast response: %w", err)
		}

		// Get the current values for the day (using current hour as index)
		currentHour := time.Now().Hour()
		if currentHour >= len(forecastData.Hourly.Relativehumidity2M) {
			return fmt.Errorf("forecast data not available for the current hour")
		}

		var FetchedHumidityCurrent = forecastData.Hourly.Relativehumidity2M[currentHour]
		var FetchedRealFeelCurrent = forecastData.Hourly.ApparentTemperature[currentHour]
		var FetchedSurfacePressureCurrent = forecastData.Hourly.SurfacePressure[currentHour]
		var FetchedSealevelPressureCurrent = forecastData.Hourly.PressureMsl[currentHour]

		airqualityUrl := AirQualityUrl + "?timezone=auto" + "&latitude=" + fmt.Sprintf("%.4f", FetchedLatitude) + "&longitude=" + fmt.Sprintf("%.4f", FetchedLongitude) + "&hourly=uv_index"
		resp, err = http.Get(airqualityUrl)
		if err != nil {
			return fmt.Errorf("failed to reach air quality API: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("air quality API returned status: %s", resp.Status)
		}

		var airqualityData AirQualityResponse
		if err := json.NewDecoder(resp.Body).Decode(&airqualityData); err != nil {
			return fmt.Errorf("failed to decode air quality response: %w", err)
		}

		// Get maximum UV index
		var FetchedUVIndexMax float64
		for _, uv := range airqualityData.Hourly.UvIndex {
			if uv > FetchedUVIndexMax {
				FetchedUVIndexMax = uv
			}
		}

		if cmd.Flag("raw").Value.String() == "true" {
			jsn, err := json.Marshal(result)
			if err != nil {
				return fmt.Errorf("failed to marshal city data: %w", err)
			}
			os.Stdout.Write(jsn)
			fmt.Println()

			jsn, err = json.Marshal(forecastData)
			if err != nil {
				return fmt.Errorf("failed to marshal forecast data: %w", err)
			}
			os.Stdout.Write(jsn)
			fmt.Println()

			jsn, err = json.Marshal(airqualityData)
			if err != nil {
				return fmt.Errorf("failed to marshal air quality data: %w", err)
			}
			os.Stdout.Write(jsn)
		} else {
			printer(FetchedCityName,
				FetchedCountryName,
				FetchedLatitude,
				FetchedLongitude,
				FetchedTimezone,
				FetchedPopulationInt,
				forecastData.CurrentWeather.Temperature,
				forecastData.CurrentWeather.Windspeed,
				forecastData.CurrentWeather.Winddirection,
				translateweathercode(fmt.Sprintf("%v", forecastData.CurrentWeather.Weathercode)),
				FetchedHumidityCurrent,
				FetchedRealFeelCurrent,
				FetchedSurfacePressureCurrent,
				FetchedSealevelPressureCurrent,
				FetchedUVIndexMax,
			)
		}
		return nil
	},
}

type GeocodingResponse struct {
	Results []struct {
		Name       string  `json:"name"`
		Country    string  `json:"country"`
		Latitude   float64 `json:"latitude"`
		Longitude  float64 `json:"longitude"`
		Timezone   string  `json:"timezone"`
		Population float64 `json:"population"`
	} `json:"results"`
}

type ForecastResponse struct {
	CurrentWeather struct {
		Temperature   float64 `json:"temperature"`
		Windspeed     float64 `json:"windspeed"`
		Winddirection float64 `json:"winddirection"`
		Weathercode   float64 `json:"weathercode"`
	} `json:"current_weather"`
	Hourly struct {
		Relativehumidity2M  []float64 `json:"relativehumidity_2m"`
		ApparentTemperature []float64 `json:"apparent_temperature"`
		SurfacePressure     []float64 `json:"surface_pressure"`
		PressureMsl         []float64 `json:"pressure_msl"`
	} `json:"hourly"`
}

type AirQualityResponse struct {
	Hourly struct {
		UvIndex []float64 `json:"uv_index"`
	} `json:"hourly"`
}


func printer(Name interface{},
	Country interface{},
	Latitude interface{},
	Longitude interface{},
	Timezone interface{},
	PopulationInt int64,
	Temperature interface{},
	WindSpeed interface{},
	WindDirection interface{},
	WeatherCode string,
	HumidityCurrent float64,
	RealFeelCurrent float64,
	SurfacePressureCurrent float64,
	SealevelPressureCurrent float64,
	UVIndexMax float64) {
	fmt.Printf("City/Country: %s/%s\n", Name, Country)
	fmt.Printf("Latitude: %f\n", Latitude)
	fmt.Printf("Longitude: %f\n", Longitude)
	fmt.Printf("Timezone: %s\n", Timezone)
	fmt.Printf("Population: %s\n", humanize.Comma(PopulationInt))
	fmt.Println("\nWeather Info:")
	fmt.Printf("	Temperature: %.1f°C\n", Temperature)
	fmt.Printf("	Wind Direction: %.0f°\n", WindDirection)
	fmt.Printf("	Wind Speed: %.1f Km/h\n", WindSpeed)
	fmt.Printf("	Weather Condition: %s\n", WeatherCode)
	fmt.Printf("	Humidity: %.2f%%\n", HumidityCurrent)
	fmt.Printf("	Real Feel: %.1f°C\n", RealFeelCurrent)
	fmt.Printf("	Surface Pressure: %.2f hPa\n", SurfacePressureCurrent)
	fmt.Printf("	Sealevel Pressure: %.2f hPa\n", SealevelPressureCurrent)
	fmt.Printf("	UV Index: %v\n", math.Round(UVIndexMax))
}

func translateweathercode(code string) string {
	switch code {
	case "0":
		return "Clear Sky"
	case "1":
		return "Mainly Clear"
	case "2":
		return "Partly Cloudy"
	case "3":
		return "Overcast"
	case "45":
		return "Fog"
	case "48":
		return "Depositing Rime Fog"
	case "51":
		return "Light Drizzle"
	case "53":
		return "Moderate Drizzle"
	case "55":
		return "Dense Drizzle"
	case "56":
		return "Light Freezing Drizzle"
	case "57":
		return "Dense Freezing Drizzle"
	case "61":
		return "Slight Rain"
	case "63":
		return "Moderate Rain"
	case "65":
		return "Heavy Rain"
	case "66":
		return "Light Freezing Rain"
	case "67":
		return "Heavy Freezing Rain"
	case "71":
		return "Slight Snow Fall"
	case "73":
		return "Moderate Snow Fall"
	case "75":
		return "Heavy Snow Fall"
	case "77":
		return "Snow Grains"
	case "80":
		return "Slight Rain Showers"
	case "81":
		return "Moderate Rain Showers"
	case "82":
		return "Violent Rain Showers"
	case "85":
		return "Slight Snow Showers"
	case "86":
		return "Heavy Snow Showers"
	case "95":
		return "Thunderstorm"
	case "96":
		return "Thunderstorm With Light Hail"
	case "99":
		return "Thunderstorm With Heavy Hail"

	default:
		return "Unknown"
	}
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolP("raw", "r", false, "Get raw data")
}
