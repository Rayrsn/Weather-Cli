package ui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	TitleBackground string
	TitleForeground string
	LabelForeground string
	ValueForeground string
	BorderColor     string
}

var (
	VibrantTheme = Theme{
		TitleBackground: "#7D56F4",
		TitleForeground: "#FAFAFA",
		LabelForeground: "#04B575",
		ValueForeground: "#FAFAFA",
		BorderColor:     "#874BFD",
	}

	DarkTheme = Theme{
		TitleBackground: "#333333",
		TitleForeground: "#EEEEEE",
		LabelForeground: "#999999",
		ValueForeground: "#FFFFFF",
		BorderColor:     "#444444",
	}

	LightTheme = Theme{
		TitleBackground: "#DDDDDD",
		TitleForeground: "#222222",
		LabelForeground: "#555555",
		ValueForeground: "#000000",
		BorderColor:     "#AAAAAA",
	}
)

func GetTheme(name string) Theme {
	switch name {
	case "dark":
		return DarkTheme
	case "light":
		return LightTheme
	default:
		return VibrantTheme
	}
}

func GetStyles(theme Theme) (title, label, value, container lipgloss.Style) {
	title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(theme.TitleForeground)).
		Background(lipgloss.Color(theme.TitleBackground)).
		Padding(0, 1).
		MarginBottom(1)

	label = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.LabelForeground)).
		Bold(true)

	value = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.ValueForeground))

	container = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.BorderColor)).
		Padding(1).
		Margin(1)

	return
}
