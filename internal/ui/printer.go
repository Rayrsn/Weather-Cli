package ui

import (
	"fmt"
	"math"
	"time"

	"github.com/Rayrsn/weather-Cli/internal/api"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

func TranslateWeatherCode(code string) string {
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

func Printer(name string,
	country string,
	latitude float64,
	longitude float64,
	timezone string,
	population int64,
	forecastData api.ForecastResponse,
	weatherCondition string,
	humidity float64,
	realFeel float64,
	surfacePressure float64,
	sealevelPressure float64,
	uvIndexMax float64,
	styled bool,
	isImperial bool,
	showForecast bool,
	showHourly bool,
	themeName string) {

	tempUnit := "°C"
	windUnit := "Km/h"
	if isImperial {
		tempUnit = "°F"
		windUnit = "mph"
	}

	weatherCode := fmt.Sprintf("%.0f", forecastData.CurrentWeather.Weathercode)

	if !styled {
		fmt.Printf("City/Country: %s/%s\n", name, country)
		fmt.Printf("Latitude: %f\n", latitude)
		fmt.Printf("Longitude: %f\n", longitude)
		fmt.Printf("Timezone: %s\n", timezone)
		fmt.Printf("Population: %s\n", humanize.Comma(population))
		fmt.Println("\nWeather Info:")
		fmt.Printf("	Temperature: %.1f%s\n", forecastData.CurrentWeather.Temperature, tempUnit)
		fmt.Printf("	Wind Direction: %.0f°\n", forecastData.CurrentWeather.Winddirection)
		fmt.Printf("	Wind Speed: %.1f %s\n", forecastData.CurrentWeather.Windspeed, windUnit)
		fmt.Printf("	Weather Condition: %s %s\n", weatherCondition, GetSmallIcon(weatherCode))
		fmt.Printf("	Humidity: %.2f%%\n", humidity)
		fmt.Printf("	Real Feel: %.1f%s\n", realFeel, tempUnit)
		fmt.Printf("	Surface Pressure: %.2f hPa\n", surfacePressure)
		fmt.Printf("	Sealevel Pressure: %.2f hPa\n", sealevelPressure)
		fmt.Printf("	UV Index: %.0f\n", math.Round(uvIndexMax))

		if showHourly && len(forecastData.Hourly.Time) > 0 {
			fmt.Println("\n24-Hour Forecast:")
			currentHour := time.Now().Hour()
			for i := currentHour; i < currentHour+24 && i < len(forecastData.Hourly.Time); i++ {
				t, _ := time.Parse("2006-01-02T15:04", forecastData.Hourly.Time[i])
				fmt.Printf("	%s: %.1f%s - %s %s\n",
					t.Format("15:04"),
					forecastData.Hourly.Temperature2M[i], tempUnit,
					TranslateWeatherCode(fmt.Sprintf("%.0f", forecastData.Hourly.Weathercode[i])),
					GetSmallIcon(fmt.Sprintf("%.0f", forecastData.Hourly.Weathercode[i])),
				)
			}
		}

		if showForecast && len(forecastData.Daily.Time) > 0 {
			fmt.Println("\n7-Day Forecast:")
			for i := 0; i < len(forecastData.Daily.Time); i++ {
				fmt.Printf("	%s: %.1f%s / %.1f%s - %s %s\n",
					forecastData.Daily.Time[i],
					forecastData.Daily.Temperature2MMax[i], tempUnit,
					forecastData.Daily.Temperature2MMin[i], tempUnit,
					TranslateWeatherCode(fmt.Sprintf("%.0f", forecastData.Daily.Weathercode[i])),
					GetSmallIcon(fmt.Sprintf("%.0f", forecastData.Daily.Weathercode[i])),
				)
			}
		}
		return
	}

	theme := GetTheme(themeName)
	titleStyle, labelStyle, valueStyle, containerStyle := GetStyles(theme)

	rowStyle := lipgloss.NewStyle().Width(60)

	renderRow := func(label, value string) string {
		l := labelStyle.Width(30).Render(label)
		v := valueStyle.Render(value)
		return rowStyle.Render(l + v)
	}

	icon := GetWeatherIcon(weatherCode)
	iconStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.LabelForeground)).MarginRight(2)
	
	locationInfo := lipgloss.JoinVertical(lipgloss.Left,
		renderRow("📍 Location:", fmt.Sprintf("%v, %v", name, country)),
		renderRow("🌐 Latitude:", fmt.Sprintf("%v", latitude)),
		renderRow("🌐 Longitude:", fmt.Sprintf("%v", longitude)),
		renderRow("🕒 Timezone:", fmt.Sprintf("%s", timezone)),
		renderRow("👥 Population:", humanize.Comma(population)),
	)

	weatherInfo := lipgloss.JoinVertical(lipgloss.Left,
		renderRow("🔥 Temp:", fmt.Sprintf("%.1f%s", forecastData.CurrentWeather.Temperature, tempUnit)),
		renderRow("💨 Wind:", fmt.Sprintf("%.1f %s (%.0f°)", forecastData.CurrentWeather.Windspeed, windUnit, forecastData.CurrentWeather.Winddirection)),
		renderRow("🌈 Condition:", weatherCondition),
		renderRow("💧 Humidity:", fmt.Sprintf("%.1f%%", humidity)),
		renderRow("🔥 Feels Like:", fmt.Sprintf("%.1f%s", realFeel, tempUnit)),
		renderRow("🧭 Surface Pressure:", fmt.Sprintf("%.2f hPa", surfacePressure)),
		renderRow("🌊 Sealevel Pressure:", fmt.Sprintf("%.2f hPa", sealevelPressure)),
		renderRow("🌞 UV Index:", fmt.Sprintf("%.0f", math.Round(uvIndexMax))),
	)

	topContent := lipgloss.JoinHorizontal(lipgloss.Top,
		iconStyle.Render(icon),
		weatherInfo,
	)

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("WEATHER REPORT"),
		locationInfo,
		"",
		titleStyle.Render("CURRENT CONDITIONS"),
		topContent,
	)

	if showHourly && len(forecastData.Hourly.Time) > 0 {
		hourlyTitle := titleStyle.MarginTop(1).Render("24-HOUR FORECAST")
		var hourlyRows []string
		currentHour := time.Now().Hour()
		for i := currentHour; i < currentHour+24 && i < len(forecastData.Hourly.Time); i++ {
			t, _ := time.Parse("2006-01-02T15:04", forecastData.Hourly.Time[i])
			timeStr := t.Format("15:04")
			
			condition := TranslateWeatherCode(fmt.Sprintf("%.0f", forecastData.Hourly.Weathercode[i]))
			temp := fmt.Sprintf("%.1f%s", forecastData.Hourly.Temperature2M[i], tempUnit)
			icon := GetSmallIcon(fmt.Sprintf("%.0f", forecastData.Hourly.Weathercode[i]))
			
			row := renderRow("🕒  "+timeStr+":", fmt.Sprintf("%-10s %s %s", temp, icon, condition))
			hourlyRows = append(hourlyRows, row)
		}
		content = lipgloss.JoinVertical(lipgloss.Left,
			content,
			hourlyTitle,
			lipgloss.JoinVertical(lipgloss.Left, hourlyRows...),
		)
	}

	if showForecast && len(forecastData.Daily.Time) > 0 {
		forecastTitle := titleStyle.MarginTop(1).Render("7-DAY FORECAST")
		var forecastRows []string
		for i := 0; i < len(forecastData.Daily.Time); i++ {
			date, _ := time.Parse("2006-01-02", forecastData.Daily.Time[i])
			dateStr := date.Format("Mon, Jan 02")

			condition := TranslateWeatherCode(fmt.Sprintf("%.0f", forecastData.Daily.Weathercode[i]))
			tempRange := fmt.Sprintf("%.1f/%.1f%s", forecastData.Daily.Temperature2MMax[i], forecastData.Daily.Temperature2MMin[i], tempUnit)

			row := renderRow("🗓️  "+dateStr+":", fmt.Sprintf("%-15s %s", tempRange, condition))
			forecastRows = append(forecastRows, row)
		}
		content = lipgloss.JoinVertical(lipgloss.Left,
			content,
			forecastTitle,
			lipgloss.JoinVertical(lipgloss.Left, forecastRows...),
		)
	}

	fmt.Println(containerStyle.Render(content))
}
