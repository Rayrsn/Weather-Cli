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
	Use:   "get [city]",
	Short: "Gets the weather for a city",
	Long:  `Gets the weather info for a city. (Can be used with --raw to get a json response)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient()
		var cityName string
		if len(args) == 0 {
			cityName = viper.GetString("default_city")
			if cityName == "" {
				// Try auto-location
				var err error
				cityName, err = client.GetAutoLocation()
				if err != nil {
					return fmt.Errorf("please enter a city name or set a default_city in your config (auto-location failed: %v)", err)
				}
			}
		} else {
			cityName = args[0]
		}

		raw, _ := cmd.Flags().GetBool("raw")
		noCache, _ := cmd.Flags().GetBool("no-cache")

		// Get units and forecast preference early for cache key
		units, _ := cmd.Flags().GetString("units")
		if units == "metric" && viper.GetString("units") != "" {
			units = viper.GetString("units")
		}
		isImperial := units == "imperial"

		showForecast, _ := cmd.Flags().GetBool("forecast")
		if !showForecast {
			showForecast = viper.GetBool("forecast")
		}

		// Try cache
		var cacheStore *cache.Cache
		var cacheKey string
		if !noCache {
			var err error
			cacheStore, err = cache.NewCache()
			if err == nil {
				cacheKey = cache.GenerateKey(cityName, isImperial, showForecast)
				if cachedData, ok := cacheStore.Get(cacheKey, 15*time.Minute); ok {
					if raw {
						jsn, _ := json.Marshal(cachedData)
						os.Stdout.Write(jsn)
						fmt.Println()
					} else {
						displayResults(cachedData.Location, &cachedData.Forecast, &cachedData.AirQuality, isImperial, showForecast, cmd)
					}
					return nil
				}
			}
		}

		var cityinfoData *api.GeocodingResponse
		if !raw {
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

			if _, err := p.Run(); err != nil {
				return fmt.Errorf("loading error: %w", err)
			}
			if cityinfoData == nil {
				return fmt.Errorf("city lookup failed for: %s", cityName)
			}
		} else {
			var err error
			cityinfoData, err = client.GetCityInfo(cityName)
			if err != nil {
				return err
			}
		}

		if len(cityinfoData.Results) == 0 {
			return fmt.Errorf("city not found: %s. Please check the spelling or try a more specific name (e.g., 'London, UK')", cityName)
		}

		noList, _ := cmd.Flags().GetBool("no-list")
		interactive := !noList

		// If flag wasn't used, check config
		if !cmd.Flags().Changed("no-list") && viper.IsSet("interactive") {
			interactive = viper.GetBool("interactive")
		}

		var result *api.GeocodingResult
		if len(cityinfoData.Results) > 1 && !raw && interactive {
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

		if !raw {
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

			if _, err := p.Run(); err != nil {
				return fmt.Errorf("loading error: %w", err)
			}
			if forecastData == nil || airqualityData == nil {
				return fmt.Errorf("failed to fetch weather or air quality data")
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

			for err := range errChan {
				if err != nil {
					return err
				}
			}
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

		if raw {
			jsn, err := json.Marshal(combined)
			if err != nil {
				return fmt.Errorf("failed to marshal combined data: %w", err)
			}
			os.Stdout.Write(jsn)
			fmt.Println()
		} else {
			displayResults(combined.Location, &combined.Forecast, &combined.AirQuality, isImperial, showForecast, cmd)
		}
		return nil
	},
}

func displayResults(location api.GeocodingResult, forecast *api.ForecastResponse, airquality *api.AirQualityResponse, isImperial bool, showForecast bool, cmd *cobra.Command) {
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
		theme,
	)
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolP("raw", "r", false, "Get raw data")
	getCmd.Flags().Bool("no-style", false, "Disable styled output")
	getCmd.Flags().StringP("units", "u", "metric", "Units to use (metric or imperial)")
	getCmd.Flags().BoolP("forecast", "f", false, "Show 7-day forecast")
	getCmd.Flags().BoolP("no-list", "l", false, "Disable interactive city selection")
	getCmd.Flags().Bool("no-cache", false, "Disable local caching")
	getCmd.Flags().StringP("theme", "t", "vibrant", "Theme to use (vibrant, dark, light)")
}
