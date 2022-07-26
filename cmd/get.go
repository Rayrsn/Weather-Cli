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
	Long:  `Gets the weather info for a city. (Can be used with --raw to get a json respons)`,
	Run: func(cmd *cobra.Command, args []string) {
		var CityName = args[0]
		if cmd.Flag("raw").Value.String() == "false" {
			fmt.Printf("Searching for city %s ...\n", strings.ToUpper(CityName[:1])+CityName[1:])
		}
		url := MainUrl + "?name=" + CityName + "&count=1"
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		var data map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&data)
		var FetchedCityName = data["results"].([]interface{})[0].(map[string]interface{})["name"]
		var FetchedCountryName = data["results"].([]interface{})[0].(map[string]interface{})["country"]
		var FetchedLatitude = data["results"].([]interface{})[0].(map[string]interface{})["latitude"]
		var FetchedLongitude = data["results"].([]interface{})[0].(map[string]interface{})["longitude"]
		var FetchedTimezone = data["results"].([]interface{})[0].(map[string]interface{})["timezone"]
		var FetchedPopulation = data["results"].([]interface{})[0].(map[string]interface{})["population"]
		var FetchedPopulationString = strconv.Itoa(int(FetchedPopulation.(float64)))
		var FetchedPopulationFloat, _ = strconv.ParseFloat(FetchedPopulationString, 64)
		var FetchedPopulationInt = int64(FetchedPopulationFloat)
		if cmd.Flag("raw").Value.String() == "true" {
			json, err := json.Marshal(data)
			if err != nil {
				log.Fatalln(err)
			}
			os.Stdout.Write(json)
		} else {
			printer(FetchedCityName, FetchedCountryName, FetchedLatitude, FetchedLongitude, FetchedTimezone, FetchedPopulation, FetchedPopulationInt)
		}
	},
}

func printer(Name interface{}, Country interface{}, Latitude interface{}, Longitude interface{}, Timezone interface{}, PopulationFloat interface{}, PopulationInt int64) {
	fmt.Printf("City/Country: %s/%s\n", Name, Country)
	fmt.Printf("Latitude: %f\n", Latitude)
	fmt.Printf("Longitude: %f\n", Longitude)
	fmt.Printf("Timezone: %s\n", Timezone)
	fmt.Printf("Population: %s (%v)\n", humanize.Comma(PopulationInt), PopulationFloat)
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolP("raw", "r", false, "Get raw data")
}
