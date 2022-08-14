/*
Copyright © 2022 Rayr https://rayr.ml/LinkInBio/
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

var GeocodingUrl = "https://geocoding-api.open-meteo.com/v1/search"
var ForecastUrl = "https://api.open-meteo.com/v1/forecast"

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets the weather for a city",
	Long:  `Gets the weather info for a city. (Can be used with --raw to get a json response)`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please enter a city name")
			os.Exit(1)
		}
		var CityNameFormatted = strings.Replace(args[0], " ", "%20", -1)
		var CityName = args[0]

		if cmd.Flag("raw").Value.String() == "false" {
			fmt.Printf("Searching for city %s...\n\n", strings.ToUpper(CityName[:1])+CityName[1:])
		}
		cityinfoUrl := GeocodingUrl + "?name=" + CityNameFormatted + "&count=1"

		resp, err := http.Get(cityinfoUrl)
		if err != nil {
			log.Fatalln(err)
		}
		var cityinfoData map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&cityinfoData)
		var FetchedCityName = cityinfoData["results"].([]interface{})[0].(map[string]interface{})["name"]
		var FetchedCountryName = cityinfoData["results"].([]interface{})[0].(map[string]interface{})["country"]
		var FetchedLatitude = cityinfoData["results"].([]interface{})[0].(map[string]interface{})["latitude"]
		var FetchedLongitude = cityinfoData["results"].([]interface{})[0].(map[string]interface{})["longitude"]
		var FetchedTimezone = cityinfoData["results"].([]interface{})[0].(map[string]interface{})["timezone"]
		var FetchedPopulation = cityinfoData["results"].([]interface{})[0].(map[string]interface{})["population"]
		var FetchedPopulationString string
		var FetchedPopulationFloat float64
		var FetchedPopulationInt int64

		if FetchedPopulation == nil {
			FetchedPopulationInt = int64(0)
		} else {
			FetchedPopulationString = strconv.Itoa(int(FetchedPopulation.(float64)))
			FetchedPopulationFloat, _ = strconv.ParseFloat(FetchedPopulationString, 64)
			FetchedPopulationInt = int64(FetchedPopulationFloat)
		}

		forecastUrl := ForecastUrl + "?latitude=" + fmt.Sprintf("%.4f", FetchedLatitude) + "&longitude=" + fmt.Sprintf("%.4f", FetchedLongitude) + "&current_weather=true" + "&hourly=relativehumidity_2m,apparent_temperature,surface_pressure"
		resp, err = http.Get(forecastUrl)
		if err != nil {
			log.Fatalln(err)
		}
		var forecastData map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&forecastData)
		var FetchedTemperature = forecastData["current_weather"].(map[string]interface{})["temperature"]
		var FetchedWindSpeed = forecastData["current_weather"].(map[string]interface{})["windspeed"]
		var FetchedWindDirection = forecastData["current_weather"].(map[string]interface{})["winddirection"]
		var FetchedWeathercode = forecastData["current_weather"].(map[string]interface{})["weathercode"]
		var FetchedHumidity = forecastData["hourly"].(map[string]interface{})["relativehumidity_2m"]
		var FetchedRealFeel = forecastData["hourly"].(map[string]interface{})["apparent_temperature"]
		var FetchedSurfacePressure = forecastData["hourly"].(map[string]interface{})["surface_pressure"]

		// Get the average humidity for the day
		var FetchedHumidityAverage float64
		for i := 0; i < 24; i++ {
			FetchedHumidityAverage += FetchedHumidity.([]interface{})[i].(float64)
		}
		FetchedHumidityAverage = FetchedHumidityAverage / 23

		// Add FetchedHumidityAverage to forecastData
		forecastData["hourly"].(map[string]interface{})["relativehumidity_2m"] = FetchedHumidityAverage

		// Get the average feels-like temperature for the day
		var FetchedRealFeelAverage float64
		for i := 0; i < 24; i++ {
			FetchedRealFeelAverage += FetchedRealFeel.([]interface{})[i].(float64)
		}
		FetchedRealFeelAverage = FetchedRealFeelAverage / 23

		// Add FetchedRealFeelAverage to forecastData
		forecastData["hourly"].(map[string]interface{})["apparent_temperature"] = FetchedRealFeelAverage

		// Get the average surface pressure for the day
		var FetchedSurfacePressureAverage float64
		for i := 0; i < 24; i++ {
			FetchedSurfacePressureAverage += FetchedSurfacePressure.([]interface{})[i].(float64)
		}
		FetchedSurfacePressureAverage = FetchedSurfacePressureAverage / 23

		// Add FetchedSurfacePressureAverage to forecastData
		forecastData["hourly"].(map[string]interface{})["surface_pressure"] = FetchedSurfacePressureAverage

		// Remove time key from hourly
		delete(forecastData["hourly"].(map[string]interface{}), "time")

		if cmd.Flag("raw").Value.String() == "true" {
			// only print the first entry under "results"
			jsn, err := json.Marshal(cityinfoData["results"].([]interface{})[0])
			if err != nil {
				log.Fatalln(err)
			}

			os.Stdout.Write(jsn)
			fmt.Println()
			jsn, err = json.Marshal(forecastData)
			if err != nil {
				log.Fatalln(err)
			}
			os.Stdout.Write(jsn)
		} else {
			printer(FetchedCityName,
				FetchedCountryName,
				FetchedLatitude,
				FetchedLongitude,
				FetchedTimezone,
				FetchedPopulation,
				FetchedPopulationInt,
				FetchedTemperature,
				FetchedWindSpeed,
				FetchedWindDirection,
				translateweathercode(fmt.Sprintf("%v", FetchedWeathercode)),
				FetchedHumidityAverage,
				FetchedRealFeelAverage,
				FetchedSurfacePressureAverage,
			)
		}
	},
}

func printer(Name interface{},
	Country interface{},
	Latitude interface{},
	Longitude interface{},
	Timezone interface{},
	PopulationFloat interface{},
	PopulationInt int64,
	Temperature interface{},
	WindSpeed interface{},
	WindDirection interface{},
	WeatherCode string,
	HumidityAverage float64,
	RealFeelAverage float64,
	SurfacePressureAverage float64) {
	fmt.Printf("City/Country: %s/%s\n", Name, Country)
	fmt.Printf("Latitude: %f\n", Latitude)
	fmt.Printf("Longitude: %f\n", Longitude)
	fmt.Printf("Timezone: %s\n", Timezone)
	fmt.Printf("Population: %s (%v)\n", humanize.Comma(PopulationInt), PopulationFloat)
	fmt.Println("\nWeather Info:")
	fmt.Printf("	Temperature: %.1f°\n", Temperature)
	fmt.Printf("	Wind Direction: %.0f°\n", WindDirection)
	fmt.Printf("	Wind Speed: %.1f Km/h\n", WindSpeed)
	fmt.Printf("	Weather Condition: %s\n", WeatherCode)
	fmt.Printf("	Humidity: %.2f%%\n", HumidityAverage)
	fmt.Printf("	Real Feel: %.1f°\n", HumidityAverage)
	fmt.Printf("	Surface Pressure: %.2f hPa\n", SurfacePressureAverage)
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
