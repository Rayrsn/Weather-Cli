package ui

func GetWeatherIcon(code string) string {
	switch code {
	case "0": // Clear sky
		return "  \\ _ /  \n" +
			   " - ( ) - \n" +
			   "  /   \\  "
	case "1", "2": // Mainly clear, partly cloudy
		return "  \\ _ /  \n" +
			   " - ( )--o\n" +
			   "  / (___)"
	case "3": // Overcast
		return "    .--.  \n" +
			   " .-(    ).\n" +
			   "(___.__.___)"
	case "45", "48": // Fog
		return " _  _  _ \n" +
			   "  _  _  _\n" +
			   " _  _  _ "
	case "51", "53", "55", "56", "57": // Drizzle
		return "    .--.  \n" +
			   " .-(    ).\n" +
			   "(___.__.___)\n" +
			   "  ' ' ' '"
	case "61", "63", "65", "66", "67", "80", "81", "82": // Rain
		return "    .--.  \n" +
			   " .-(    ).\n" +
			   "(___.__.___)\n" +
			   "  / / / /"
	case "71", "73", "75", "77", "85", "86": // Snow
		return "    .--.  \n" +
			   " .-(    ).\n" +
			   "(___.__.___)\n" +
			   "  * * * *"
	case "95", "96", "99": // Thunderstorm
		return "    .--.  \n" +
			   " .-(    ).\n" +
			   "(___.__.___)\n" +
			   "  ⚡ ⚡ ⚡"
	default:
		return "    ?    \n" +
			   "   ???   \n" +
			   "    ?    "
	}
}

func GetSmallIcon(code string) string {
	switch code {
	case "0":
		return "☀️"
	case "1", "2":
		return "🌤️"
	case "3":
		return "☁️"
	case "45", "48":
		return "🌫️"
	case "51", "53", "55", "56", "57":
		return "🌦️"
	case "61", "63", "65", "66", "67", "80", "81", "82":
		return "🌧️"
	case "71", "73", "75", "77", "85", "86":
		return "❄️"
	case "95", "96", "99":
		return "⛈️"
	default:
		return "❓"
	}
}
