package cmd

import "testing"

func TestTranslateweathercode(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{"0", "Clear Sky"},
		{"1", "Mainly Clear"},
		{"2", "Partly Cloudy"},
		{"3", "Overcast"},
		{"45", "Fog"},
		{"48", "Depositing Rime Fog"},
		{"51", "Light Drizzle"},
		{"53", "Moderate Drizzle"},
		{"55", "Dense Drizzle"},
		{"56", "Light Freezing Drizzle"},
		{"57", "Dense Freezing Drizzle"},
		{"61", "Slight Rain"},
		{"63", "Moderate Rain"},
		{"65", "Heavy Rain"},
		{"66", "Light Freezing Rain"},
		{"67", "Heavy Freezing Rain"},
		{"71", "Slight Snow Fall"},
		{"73", "Moderate Snow Fall"},
		{"75", "Heavy Snow Fall"},
		{"77", "Snow Grains"},
		{"80", "Slight Rain Showers"},
		{"81", "Moderate Rain Showers"},
		{"82", "Violent Rain Showers"},
		{"85", "Slight Snow Showers"},
		{"86", "Heavy Snow Showers"},
		{"95", "Thunderstorm"},
		{"96", "Thunderstorm With Light Hail"},
		{"99", "Thunderstorm With Heavy Hail"},
		{"999", "Unknown"},
		{"", "Unknown"},
	}

	for _, tt := range tests {
		result := translateweathercode(tt.code)
		if result != tt.expected {
			t.Errorf("translateweathercode(%s) = %s; want %s", tt.code, result, tt.expected)
		}
	}
}
