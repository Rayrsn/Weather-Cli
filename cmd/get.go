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

var MainUrl = "https://geocoding-api.open-meteo.com/v1/search"
var SecondUrl = "https://api.open-meteo.com/v1/forecast"

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets the weather for a city",
	Long:  `Gets the weather info for a city. (Can be used with --raw to get a json response)`,
	Run: func(cmd *cobra.Command, args []string) {
		var CityName = args[0]
		if cmd.Flag("raw").Value.String() == "false" {
			fmt.Printf("Searching for city %s ...\n\n", strings.ToUpper(CityName[:1])+CityName[1:])
		}
		cityinfoUrl := MainUrl + "?name=" + CityName + "&count=1"

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
		var FetchedPopulationString = strconv.Itoa(int(FetchedPopulation.(float64)))
		var FetchedPopulationFloat, _ = strconv.ParseFloat(FetchedPopulationString, 64)
		var FetchedPopulationInt = int64(FetchedPopulationFloat)

		forecastUrl := SecondUrl + "?latitude=" + fmt.Sprintf("%.4f", FetchedLatitude) + "&longitude=" + fmt.Sprintf("%.4f", FetchedLongitude) + "&current_weather=true" + "&hourly=relativehumidity_2m"
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

		// Get the average humidity for the day
		var FetchedHumidityAverage float64
		for i := 0; i < 24; i++ {
			FetchedHumidityAverage += FetchedHumidity.([]interface{})[i].(float64)
		}
		FetchedHumidityAverage = FetchedHumidityAverage / 23

		if cmd.Flag("raw").Value.String() == "true" {
			jsn, err := json.Marshal(cityinfoData)
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
	HumidityAverage float64) {
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
		return "Slight Snow fall"
	case "73":
		return "Moderate Snow fall"
	case "75":
		return "Heavy Snow fall"
	case "77":
		return "Snow Grains"
	case "80":
		return "Slight Rain showers"
	case "81":
		return "Moderate Rain showers"
	case "82":
		return "Violent Rain showers"
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
