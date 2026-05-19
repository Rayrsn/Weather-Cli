/*
Copyright © 2026 Rayr https://rayrsn.me/LinkInBio
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Rayrsn/weather-Cli/internal/api"
	"github.com/Rayrsn/weather-Cli/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCmd = &cobra.Command{
	Use:   "get [city]",
	Short: "Gets the weather for a city",
	Long:  `Gets the weather info for a city. (Can be used with --raw to get a json response)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cityName string
		if len(args) == 0 {
			cityName = viper.GetString("default_city")
			if cityName == "" {
				return fmt.Errorf("please enter a city name or set a default_city in your config")
			}
		} else {
			cityName = args[0]
		}

		raw, _ := cmd.Flags().GetBool("raw")
		if !raw {
			searchStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7FF")).Bold(true)
			fmt.Printf("%s\n\n", searchStyle.Render(fmt.Sprintf("🔍 Searching for city %s...", strings.ToUpper(cityName[:1])+cityName[1:])))
		}

		client := api.NewClient()

		cityinfoData, err := client.GetCityInfo(cityName)
		if err != nil {
			return err
		}

		if len(cityinfoData.Results) == 0 {
			return fmt.Errorf("city not found: %s", cityName)
		}

		result := cityinfoData.Results[0]

		// Get units preference
		units, _ := cmd.Flags().GetString("units")
		if units == "metric" && viper.GetString("units") != "" {
			units = viper.GetString("units")
		}
		isImperial := units == "imperial"

		// Get forecast preference
		showForecast, _ := cmd.Flags().GetBool("forecast")
		if !showForecast {
			showForecast = viper.GetBool("forecast")
		}

		var forecastData *api.ForecastResponse
		var airqualityData *api.AirQualityResponse
		var wg sync.WaitGroup
		errChan := make(chan error, 2)

		wg.Add(2)
		go func() {
			defer wg.Done()
			var err error
			forecastData, err = client.GetForecast(result.Latitude, result.Longitude, isImperial, showForecast)
			if err != nil {
				errChan <- err
			}
		}()

		go func() {
			defer wg.Done()
			var err error
			airqualityData, err = client.GetAirQuality(result.Latitude, result.Longitude)
			if err != nil {
				errChan <- err
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

		humidityCurrent := forecastData.Hourly.Relativehumidity2M[currentHour]
		realFeelCurrent := forecastData.Hourly.ApparentTemperature[currentHour]
		surfacePressureCurrent := forecastData.Hourly.SurfacePressure[currentHour]
		sealevelPressureCurrent := forecastData.Hourly.PressureMsl[currentHour]

		// Get maximum UV index
		var uvIndexMax float64
		for _, uv := range airqualityData.Hourly.UvIndex {
			if uv > uvIndexMax {
				uvIndexMax = uv
			}
		}

		if raw {
			combined := api.CombinedResponse{
				Location:   result,
				Forecast:   *forecastData,
				AirQuality: *airqualityData,
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

			ui.Printer(result.Name,
				result.Country,
				result.Latitude,
				result.Longitude,
				result.Timezone,
				int64(result.Population),
				*forecastData,
				ui.TranslateWeatherCode(fmt.Sprintf("%v", forecastData.CurrentWeather.Weathercode)),
				humidityCurrent,
				realFeelCurrent,
				surfacePressureCurrent,
				sealevelPressureCurrent,
				uvIndexMax,
				!noStyle,
				isImperial,
				showForecast,
			)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolP("raw", "r", false, "Get raw data")
	getCmd.Flags().Bool("no-style", false, "Disable styled output")
	getCmd.Flags().StringP("units", "u", "metric", "Units to use (metric or imperial)")
	getCmd.Flags().BoolP("forecast", "f", false, "Show 7-day forecast")
}
