package ui

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/Rayrsn/weather-Cli/internal/api"
)

func ExportCSV(data api.CombinedResponse, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"City", "Country", "Latitude", "Longitude", "Time", "Temperature", "WeatherCode"})

	// Write current weather
	writer.Write([]string{
		data.Location.Name,
		data.Location.Country,
		fmt.Sprintf("%f", data.Location.Latitude),
		fmt.Sprintf("%f", data.Location.Longitude),
		"Current",
		fmt.Sprintf("%.1f", data.Forecast.CurrentWeather.Temperature),
		fmt.Sprintf("%.0f", data.Forecast.CurrentWeather.Weathercode),
	})

	// Write hourly
	for i := 0; i < len(data.Forecast.Hourly.Time); i++ {
		writer.Write([]string{
			data.Location.Name,
			data.Location.Country,
			fmt.Sprintf("%f", data.Location.Latitude),
			fmt.Sprintf("%f", data.Location.Longitude),
			data.Forecast.Hourly.Time[i],
			fmt.Sprintf("%.1f", data.Forecast.Hourly.Temperature2M[i]),
			fmt.Sprintf("%.0f", data.Forecast.Hourly.Weathercode[i]),
		})
	}

	return nil
}

func ExportMarkdown(data api.CombinedResponse, filePath string) error {
	content := fmt.Sprintf("# Weather Report for %s, %s\n\n", data.Location.Name, data.Location.Country)
	content += fmt.Sprintf("- **Latitude:** %f\n", data.Location.Latitude)
	content += fmt.Sprintf("- **Longitude:** %f\n", data.Location.Longitude)
	content += fmt.Sprintf("- **Population:** %.0f\n\n", data.Location.Population)

	content += "## Current Weather\n\n"
	content += fmt.Sprintf("- **Temperature:** %.1f\n", data.Forecast.CurrentWeather.Temperature)
	content += fmt.Sprintf("- **Condition:** %s\n\n", TranslateWeatherCode(fmt.Sprintf("%.0f", data.Forecast.CurrentWeather.Weathercode)))

	content += "## Hourly Forecast\n\n"
	content += "| Time | Temperature | Condition |\n"
	content += "| --- | --- | --- |\n"

	for i := 0; i < 24 && i < len(data.Forecast.Hourly.Time); i++ {
		content += fmt.Sprintf("| %s | %.1f | %s |\n",
			data.Forecast.Hourly.Time[i],
			data.Forecast.Hourly.Temperature2M[i],
			TranslateWeatherCode(fmt.Sprintf("%.0f", data.Forecast.Hourly.Weathercode[i])),
		)
	}

	return os.WriteFile(filePath, []byte(content), 0644)
}
