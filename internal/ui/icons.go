package ui

func GetWeatherIcon(code string) string {
	switch code {
	case "0": // Clear sky
		return "     \\  |  /\n" +
			   "    .-'\"\"\"'-.\n" +
			   "  -(         )-\n" +
			   "    '-.___.-'\n" +
			   "     /  |  \\"
	case "1", "2": // Mainly clear, partly cloudy
		return "   \\  |  /\n" +
			   "  .-'\"\"\"'-.\n" +
			   " -(      .---.\n" +
			   "  '-.  -(     ).\n" +
			   "   / (___._____)"
	case "3": // Overcast
		return "       .--.\n" +
			   "    .-((    )-.\n" +
			   "   (           )\n" +
			   "  (             )\n" +
			   "   '---'---'---'"
	case "45", "48": // Fog
		return "   _..---.._\n" +
			   "  '-..___..-'\n" +
			   "   _..---.._\n" +
			   "  '-..___..-'\n" +
			   "   _..---.._"
	case "51", "53", "55", "56", "57": // Drizzle
		return "       .--.\n" +
			   "    .-((    )-.\n" +
			   "   (           )\n" +
			   "  (             )\n" +
			   "   '---'---'---'\n" +
			   "     '   '   '"
	case "61", "63", "65", "66", "67", "80", "81", "82": // Rain
		return "       .--.\n" +
			   "    .-((    )-.\n" +
			   "   (           )\n" +
			   "  (             )\n" +
			   "   '---'---'---'\n" +
			   "    /  /  /  /\n" +
			   "   /  /  /  /"
	case "71", "73", "75", "77", "85", "86": // Snow
		return "       .--.\n" +
			   "    .-((    )-.\n" +
			   "   (           )\n" +
			   "  (             )\n" +
			   "   '---'---'---'\n" +
			   "     *   *   *\n" +
			   "   *   *   *"
	case "95", "96", "99": // Thunderstorm
		return "       .--.\n" +
			   "    .-((    )-.\n" +
			   "   (           )\n" +
			   "  (             )\n" +
			   "   '---'---'---'\n" +
			   "        _/_\n" +
			   "       /  /_\n" +
			   "      / /_/\n" +
			   "       /_/"
	default:
		return "     .---.\n" +
			   "    /     \\\n" +
			   "    '  .-.'\n" +
			   "      / /\n" +
			   "     (_)"
	}
}

func GetSmallIcon(code string) string {
	switch code {
	case "0":
		return "☀️"
	case "1", "2":
		return "🌤️"
	case "3":
		// Using 'Sun Behind Large Cloud' (🌥️) which is rendered consistently 
		// as an emoji and generally agrees on width across terminals.
		// The standard Cloud emoji (☁️) causes a 1-space border misalignment 
		// due to conflicts between go-runewidth and terminal renderers.
		return "🌥️"
	case "45", "48":
		return "🌫️"
	case "51", "53", "55", "56", "57":
		return "🌦️"
	case "61", "63", "65", "66", "67", "80", "81", "82":
		return "🌧️"
	case "71", "73", "75", "77", "85", "86":
		return "🌨️"
	case "95", "96", "99":
		return "⛈️"
	default:
		return "❓"
	}
}
