/*
Copyright © 2026 Rayr https://rayrsn.me/LinkInBio
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Rayrsn/weather-Cli/internal/api"
	"github.com/Rayrsn/weather-Cli/internal/cache"
	"github.com/Rayrsn/weather-Cli/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCmd = &cobra.Command{
	Use:   "get [city...]",
	Short: "Gets the weather for one or more cities",
	Long:  `Gets the weather info for one or more cities. (Can be used with --raw to get a json response)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient()
		var cities []string

		if len(args) == 0 {
			defaultCity := viper.GetString("default_city")
			if defaultCity == "" {
				// Try auto-location
				autoCity, err := client.GetAutoLocation()
				if err != nil {
					return fmt.Errorf("please enter a city name or set a default_city in your config (auto-location failed: %v)", err)
				}
				cities = append(cities, autoCity)
			} else {
				cities = append(cities, defaultCity)
			}
		} else {
			cities = args
		}

		raw, _ := cmd.Flags().GetBool("raw")
		noCache, _ := cmd.Flags().GetBool("no-cache")
		output, _ := cmd.Flags().GetString("output")
		outputFile, _ := cmd.Flags().GetString("output-file")

		// Get units and forecast preference early
		units, _ := cmd.Flags().GetString("units")
		if units == "metric" && viper.GetString("units") != "" {
			units = viper.GetString("units")
		}
		isImperial := units == "imperial"

		showForecast, _ := cmd.Flags().GetBool("forecast")
		if !showForecast {
			showForecast = viper.GetBool("forecast")
		}

		showHourly, _ := cmd.Flags().GetBool("hourly")
		if !showHourly {
			showHourly = viper.GetBool("hourly")
		}

		cacheStore, _ := cache.NewCache()

		for _, cityName := range cities {
			if len(cities) > 1 && !raw && output == "" {
				fmt.Printf("\n--- Weather for %s ---\n", cityName)
			}

			// Try cache
			cacheKey := cache.GenerateKey(cityName, isImperial, showForecast)
			if !noCache && cacheStore != nil {
				if cachedData, ok := cacheStore.Get(cacheKey, 15*time.Minute); ok {
					handleOutput(*cachedData, raw, output, outputFile, isImperial, showForecast, showHourly, cmd)
					continue
				}
			}

			var cityinfoData *api.GeocodingResponse
			if !raw && output == "" {
				loader := ui.NewLoadingModel(fmt.Sprintf("Searching for city %s", cityName))
				p := tea.NewProgram(loader)

				go func() {
					var err error
					cityinfoData, err = client.GetCityInfo(cityName)
					if err != nil {
						p.Send(err)
						return
					}
					p.Send(true)
				}()

				m, err := p.Run()
				if err != nil {
					return fmt.Errorf("loading error: %w", err)
				}
				if loaderModel, ok := m.(ui.LoadingModel); ok && loaderModel.Err() != nil {
					fmt.Printf("Error searching for %s: %v\n", cityName, loaderModel.Err())
					continue
				}
			} else {
				var err error
				cityinfoData, err = client.GetCityInfo(cityName)
				if err != nil {
					fmt.Printf("Error searching for %s: %v\n", cityName, err)
					continue
				}
			}

			if cityinfoData == nil || len(cityinfoData.Results) == 0 {
				fmt.Printf("Error: city not found: %s\n", cityName)
				continue
			}

			noList, _ := cmd.Flags().GetBool("no-list")
			interactive := !noList && len(cities) == 1 && output == "" && !raw

			if !cmd.Flags().Changed("no-list") && viper.IsSet("interactive") {
				interactive = viper.GetBool("interactive") && len(cities) == 1 && output == "" && !raw
			}

			var result *api.GeocodingResult
			if len(cityinfoData.Results) > 1 && interactive {
				var err error
				result, err = ui.SelectCity(cityinfoData.Results)
				if err != nil {
					return err
				}
			} else {
				result = &cityinfoData.Results[0]
			}

			var forecastData *api.ForecastResponse
			var airqualityData *api.AirQualityResponse

			if !raw && output == "" {
				loader := ui.NewLoadingModel(fmt.Sprintf("Fetching weather data for %s", result.Name))
				p := tea.NewProgram(loader)

				go func() {
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
							p.Send(err)
							return
						}
					}
					p.Send(true)
				}()

				m, err := p.Run()
				if err != nil {
					return fmt.Errorf("loading error: %w", err)
				}
				if loaderModel, ok := m.(ui.LoadingModel); ok && loaderModel.Err() != nil {
					fmt.Printf("Error fetching weather for %s: %v\n", result.Name, loaderModel.Err())
					continue
				}
			} else {
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

				fetchFailed := false
				for err := range errChan {
					if err != nil {
						fmt.Printf("Error fetching data: %v\n", err)
						fetchFailed = true
						break
					}
				}
				if fetchFailed {
					continue
				}
			}

			if forecastData == nil || airqualityData == nil {
				continue
			}

			combined := api.CombinedResponse{
				Location:   *result,
				Forecast:   *forecastData,
				AirQuality: *airqualityData,
			}

			// Save to cache
			if cacheStore != nil {
				cacheStore.Set(cacheKey, combined)
			}

			handleOutput(combined, raw, output, outputFile, isImperial, showForecast, showHourly, cmd)
		}
		return nil
	},
}

func handleOutput(data api.CombinedResponse, raw bool, output, outputFile string, isImperial, showForecast, showHourly bool, cmd *cobra.Command) {
	if raw || output == "json" {
		jsn, _ := json.Marshal(data)
		if outputFile != "" {
			os.WriteFile(outputFile, jsn, 0644)
			fmt.Printf("Exported JSON to %s\n", outputFile)
		} else {
			os.Stdout.Write(jsn)
			fmt.Println()
		}
		return
	}

	if output == "csv" {
		if outputFile == "" {
			outputFile = fmt.Sprintf("weather_%s.csv", data.Location.Name)
		}
		if err := ui.ExportCSV(data, outputFile); err != nil {
			fmt.Printf("Error exporting CSV: %v\n", err)
		} else {
			fmt.Printf("Exported CSV to %s\n", outputFile)
		}
		return
	}

	if output == "markdown" || output == "md" {
		if outputFile == "" {
			outputFile = fmt.Sprintf("weather_%s.md", data.Location.Name)
		}
		if err := ui.ExportMarkdown(data, outputFile); err != nil {
			fmt.Printf("Error exporting Markdown: %v\n", err)
		} else {
			fmt.Printf("Exported Markdown to %s\n", outputFile)
		}
		return
	}

	displayResults(data.Location, &data.Forecast, &data.AirQuality, isImperial, showForecast, showHourly, cmd)
}

func displayResults(location api.GeocodingResult, forecast *api.ForecastResponse, airquality *api.AirQualityResponse, isImperial bool, showForecast bool, showHourly bool, cmd *cobra.Command) {
	currentHour := time.Now().Hour()
	if currentHour >= len(forecast.Hourly.Relativehumidity2M) {
		currentHour = len(forecast.Hourly.Relativehumidity2M) - 1
	}

	humidityCurrent := forecast.Hourly.Relativehumidity2M[currentHour]
	realFeelCurrent := forecast.Hourly.ApparentTemperature[currentHour]
	surfacePressureCurrent := forecast.Hourly.SurfacePressure[currentHour]
	sealevelPressureCurrent := forecast.Hourly.PressureMsl[currentHour]

	var uvIndexMax float64
	for _, uv := range airquality.Hourly.UvIndex {
		if uv > uvIndexMax {
			uvIndexMax = uv
		}
	}

	// Simple alert based on weather code
	if forecast.CurrentWeather.Weathercode >= 95 {
		fmt.Printf("\n⚠️  WARNING: %s\n", ui.TranslateWeatherCode(fmt.Sprintf("%.0f", forecast.CurrentWeather.Weathercode)))
	}

	noStyle, _ := cmd.Flags().GetBool("no-style")
	if !noStyle {
		noStyle = viper.GetBool("no_style")
	}

	theme, _ := cmd.Flags().GetString("theme")
	if theme == "vibrant" && viper.GetString("theme") != "" {
		theme = viper.GetString("theme")
	}

	ui.Printer(location.Name,
		location.Country,
		location.Latitude,
		location.Longitude,
		location.Timezone,
		int64(location.Population),
		*forecast,
		ui.TranslateWeatherCode(fmt.Sprintf("%v", forecast.CurrentWeather.Weathercode)),
		humidityCurrent,
		realFeelCurrent,
		surfacePressureCurrent,
		sealevelPressureCurrent,
		uvIndexMax,
		!noStyle,
		isImperial,
		showForecast,
		showHourly,
		theme,
	)
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolP("raw", "r", false, "Get raw data")
	getCmd.Flags().Bool("no-style", false, "Disable styled output")
	getCmd.Flags().StringP("units", "u", "metric", "Units to use (metric or imperial)")
	getCmd.Flags().BoolP("forecast", "f", false, "Show 7-day forecast")
	getCmd.Flags().BoolP("hourly", "H", false, "Show 24-hour forecast")
	getCmd.Flags().BoolP("no-list", "l", false, "Disable interactive city selection")
	getCmd.Flags().Bool("no-cache", false, "Disable local caching")
	getCmd.Flags().StringP("theme", "t", "vibrant", "Theme to use (vibrant, dark, light)")
	getCmd.Flags().StringP("output", "o", "", "Output format (json, csv, markdown)")
	getCmd.Flags().String("output-file", "", "Output file path")
}
