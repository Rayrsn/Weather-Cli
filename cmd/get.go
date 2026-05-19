/*
Copyright © 2026 Rayr https://rayrsn.me/LinkInBio
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GeocodingUrl = "https://geocoding-api.open-meteo.com/v1/search"
var ForecastUrl = "https://api.open-meteo.com/v1/forecast"
var AirQualityUrl = "https://air-quality-api.open-meteo.com/v1/air-quality"

var getCmd = &cobra.Command{
	Use:   "get [city]",
	Short: "Gets the weather for a city",
	Long:  `Gets the weather info for a city. (Can be used with --raw to get a json response)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var CityName string
		if len(args) == 0 {
			CityName = viper.GetString("default_city")
			if CityName == "" {
				return fmt.Errorf("please enter a city name or set a default_city in your config")
			}
		} else {
			CityName = args[0]
		}
		var CityNameFormatted = url.QueryEscape(CityName)

		if cmd.Flag("raw").Value.String() == "false" {
			searchStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7FF")).Bold(true)
			fmt.Printf("%s\n\n", searchStyle.Render(fmt.Sprintf("🔍 Searching for city %s...", strings.ToUpper(CityName[:1])+CityName[1:])))
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

		var forecastData ForecastResponse
		var airqualityData AirQualityResponse
		var wg sync.WaitGroup
		errChan := make(chan error, 2)

		wg.Add(2)
		go func() {
			defer wg.Done()
			forecastUrl := ForecastUrl + "?timezone=auto" + "&latitude=" + fmt.Sprintf("%.4f", FetchedLatitude) + "&longitude=" + fmt.Sprintf("%.4f", FetchedLongitude) + "&current_weather=true" + "&hourly=relativehumidity_2m,apparent_temperature,surface_pressure,pressure_msl"
			resp, err := http.Get(forecastUrl)
			if err != nil {
				errChan <- fmt.Errorf("failed to reach forecast API: %w", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				errChan <- fmt.Errorf("forecast API returned status: %s", resp.Status)
				return
			}

			if err := json.NewDecoder(resp.Body).Decode(&forecastData); err != nil {
				errChan <- fmt.Errorf("failed to decode forecast response: %w", err)
				return
			}
		}()

		go func() {
			defer wg.Done()
			airqualityUrl := AirQualityUrl + "?timezone=auto" + "&latitude=" + fmt.Sprintf("%.4f", FetchedLatitude) + "&longitude=" + fmt.Sprintf("%.4f", FetchedLongitude) + "&hourly=uv_index"
			resp, err := http.Get(airqualityUrl)
			if err != nil {
				errChan <- fmt.Errorf("failed to reach air quality API: %w", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				errChan <- fmt.Errorf("air quality API returned status: %s", resp.Status)
				return
			}

			if err := json.NewDecoder(resp.Body).Decode(&airqualityData); err != nil {
				errChan <- fmt.Errorf("failed to decode air quality response: %w", err)
				return
			}
		}()

		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				return err
			}
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

		// Get maximum UV index
		var FetchedUVIndexMax float64
		for _, uv := range airqualityData.Hourly.UvIndex {
			if uv > FetchedUVIndexMax {
				FetchedUVIndexMax = uv
			}
		}

		if cmd.Flag("raw").Value.String() == "true" {
			combined := CombinedResponse{
				Location:   result,
				Forecast:   forecastData,
				AirQuality: airqualityData,
			}
			jsn, err := json.Marshal(combined)
			if err != nil {
				return fmt.Errorf("failed to marshal combined data: %w", err)
			}
			os.Stdout.Write(jsn)
			fmt.Println()
		} else {
			noStyle, _ := cmd.Flags().GetBool("no-style")
			if !noStyle {
				noStyle = viper.GetBool("no_style")
			}

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
				!noStyle,
			)
		}
		return nil
	},
}

type GeocodingResult struct {
	Name       string  `json:"name"`
	Country    string  `json:"country"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Timezone   string  `json:"timezone"`
	Population float64 `json:"population"`
}

type GeocodingResponse struct {
	Results []GeocodingResult `json:"results"`
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

type CombinedResponse struct {
	Location   GeocodingResult    `json:"location"`
	Forecast   ForecastResponse   `json:"forecast"`
	AirQuality AirQualityResponse `json:"air_quality"`
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
	UVIndexMax float64,
	styled bool) {

	if !styled {
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
		fmt.Printf("	UV Index: %.0f\n", math.Round(UVIndexMax))
		return
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA"))

	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1).
		Margin(1)

	rowStyle := lipgloss.NewStyle().Width(50)

	renderRow := func(label, value string) string {
		l := labelStyle.Width(25).Render(label)
		v := valueStyle.Render(value)
		return rowStyle.Render(l + v)
	}

	locationInfo := lipgloss.JoinVertical(lipgloss.Left,
		renderRow("📍 Location:", fmt.Sprintf("%v, %v", Name, Country)),
		renderRow("🌐 Latitude:", fmt.Sprintf("%v", Latitude)),
		renderRow("🌐 Longitude:", fmt.Sprintf("%v", Longitude)),
		renderRow("🕒 Timezone:", fmt.Sprintf("%s", Timezone)),
		renderRow("👥 Population:", humanize.Comma(PopulationInt)),
	)

	weatherInfo := lipgloss.JoinVertical(lipgloss.Left,
		renderRow("🔥 Temp:", fmt.Sprintf("%.1f°C", Temperature)),
		renderRow("💨 Wind:", fmt.Sprintf("%.1f Km/h (%.0f°)", WindSpeed, WindDirection)),
		renderRow("⛅ Condition:", WeatherCode),
		renderRow("💧 Humidity:", fmt.Sprintf("%.1f%%", HumidityCurrent)),
		renderRow("🔥 Feels Like:", fmt.Sprintf("%.1f°C", RealFeelCurrent)),
		renderRow("🧭 Surface Pressure:", fmt.Sprintf("%.2f hPa", SurfacePressureCurrent)),
		renderRow("🌊 Sealevel Pressure:", fmt.Sprintf("%.2f hPa", SealevelPressureCurrent)),
		renderRow("🌞 UV Index:", fmt.Sprintf("%.0f", math.Round(UVIndexMax))),
	)

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("WEATHER REPORT"),
		locationInfo,
		"",
		titleStyle.Render("CURRENT CONDITIONS"),
		weatherInfo,
	)

	fmt.Println(containerStyle.Render(content))
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
	getCmd.Flags().Bool("no-style", false, "Disable styled output")
}
