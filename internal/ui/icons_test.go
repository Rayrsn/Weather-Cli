package ui

import (
	"testing"
)

func TestVisualPreview(t *testing.T) {
	codes := []struct {
		code string
		name string
	}{
		{"0", "Clear sky"},
		{"1", "Mainly clear, partly cloudy"},
		{"3", "Overcast"},
		{"45", "Fog"},
		{"51", "Drizzle"},
		{"61", "Rain"},
		{"71", "Snow"},
		{"95", "Thunderstorm"},
		{"unknown", "Unknown / Default"},
	}

	for _, c := range codes {
		t.Logf("=== %s (%s) ===\n%s\n", c.name, c.code, GetWeatherIcon(c.code))
	}
}
