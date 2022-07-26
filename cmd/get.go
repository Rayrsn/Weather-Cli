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

//https://api.open-meteo.com/v1/forecast?latitude=35.7061&longitude=51.4358&hourly=temperature_2m
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets the weather for a city",
	Long:  `Gets the weather info for a city. (Can be used with --raw to get a json respons)`,
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

		forecastUrl := SecondUrl + "?latitude=" + fmt.Sprintf("%.4f", FetchedLatitude) + "&longitude=" + fmt.Sprintf("%.4f", FetchedLongitude) + "&current_weather=true"
		fmt.Println(forecastUrl)
		resp, err = http.Get(forecastUrl)
		if err != nil {
			log.Fatalln(err)
		}
		var forecastData map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&forecastData)
		fmt.Println(forecastData)
		var FetchedTemperature = forecastData["current_weather"].(map[string]interface{})["temperature"]
		var FetchedWindSpeed = forecastData["current_weather"].(map[string]interface{})["windspeed"]
		var FetchedWindDirection = forecastData["current_weather"].(map[string]interface{})["winddirection"]

		if cmd.Flag("raw").Value.String() == "true" {
			json, err := json.Marshal(cityinfoData)
			if err != nil {
				log.Fatalln(err)
			}
			os.Stdout.Write(json)
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
			)
		}
	},
}

func printer(Name interface{}, Country interface{}, Latitude interface{}, Longitude interface{}, Timezone interface{}, PopulationFloat interface{}, PopulationInt int64, Temperature interface{}, WindSpeed interface{}, WindDirection interface{}) {
	fmt.Printf("City/Country: %s/%s\n", Name, Country)
	fmt.Printf("Latitude: %f\n", Latitude)
	fmt.Printf("Longitude: %f\n", Longitude)
	fmt.Printf("Timezone: %s\n", Timezone)
	fmt.Printf("Population: %s (%v)\n", humanize.Comma(PopulationInt), PopulationFloat)
	fmt.Println("\nWeather:")
	fmt.Printf("	Temperature: %.1f°\n", Temperature)
	fmt.Printf("	Wind Direction: %.0f°\n", WindDirection)
	fmt.Printf("	Wind Speed: %.1f Km/h\n", WindSpeed)
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolP("raw", "r", false, "Get raw data")
}
