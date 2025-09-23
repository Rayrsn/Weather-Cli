package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Rayrsn/Weather-Cli/utils"
)

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")

	var weather utils.Weather
	var err error

	if city != "" {
		weather, err = utils.GetWeatherByCity(city)
	} else {
		loc, lat, lon, err2 := utils.GetCurrentLocation()
		if err2 != nil {
			http.Error(w, "Could not detect location", http.StatusInternalServerError)
			return
		}
		weather, err = utils.GetWeatherByCoords(lat, lon, loc)
	}

	if err != nil {
		http.Error(w, "Could not fetch weather", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/weather", weatherHandler)

	fmt.Println("🌐 Running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
