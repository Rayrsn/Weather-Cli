<div align="center">
  
# `🌦 Weather-Cli`
## **A modern, fast CLI tool for weather information**

<p align="center">
  <img width=40% height=40% src="https://github.com/Rayrsn/Weather-Cli/raw/main/cloud.gif?raw=true" alt="cloud">
</p>

[![Status Badge](https://github.com/Rayrsn/Weather-Cli/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/Rayrsn/Weather-Cli/actions/workflows/build.yml)

</div>

## Features 🏆
* **Modern UI:** Beautiful, styled terminal output with ASCII icons and colors using `lipgloss`.
* **Theming:** Choose between `vibrant`, `dark`, and `light` themes.
* **Extended Forecasts:** View a detailed 7-day outlook (`--forecast`) or a 24-hour breakdown (`--hourly`).
* **Interactive Selection:** Easily pick from multiple city results when names are ambiguous.
* **Auto-Location:** Automatically detect your current city using IP geolocation.
* **Caching:** Built-in 15-minute caching for near-instant repeat lookups.
* **Multi-City Support:** Query multiple cities in a single command.
* **Exporting:** Save weather reports to `JSON`, `CSV`, or `Markdown`.
* **Configuration CLI:** Manage your preferences directly through the tool.
* **Automation-Friendly:** Unified `--raw` JSON output for piping into tools like `jq`.

## Building 🔨
### Manual Build
Make sure you have [Go](https://go.dev/) installed.

```bash
go build -o weather-Cli
```

### Pre-built Binaries & Packages
* **Arch Linux:** Available in the [weather-cli](https://aur.archlinux.org/packages/weather-cli) AUR.
* **Binaries:** You can grab a binary from the [Releases](https://github.com/Rayrsn/Weather-Cli/releases).

## Usage 💻

### Get Weather for a City
```bash
./weather-Cli get London
```

### Get Multiple Cities
```bash
./weather-Cli get London Tokyo "New York"
```

### Hourly and Extended Forecasts
```bash
./weather-Cli get Paris --hourly --forecast
```

### Exporting Data
```bash
./weather-Cli get London --output markdown --output-file report.md
```

### Use Imperial Units & Different Themes
```bash
./weather-Cli get "New York" --units imperial --theme dark
```

## Examples 📸

### Standard Output
```bash
./weather-Cli get London
```
```text
  ╭───────────────────────────────────────────────────────────────────────╮ 
 │                                                                         │ 
 │  WEATHER REPORT                                                         │ 
 │                                                                         │ 
 │ 📍 Location:                  London, United Kingdom                    │ 
 │ 🌐 Latitude:                  51.50853                                  │ 
 │ 🌐 Longitude:                 -0.12574                                  │ 
 │ 🕒 Timezone:                  Europe/London                             │ 
 │ 👥 Population:                8,961,989                                 │ 
 │                                                                         │ 
 │  CURRENT CONDITIONS                                                     │ 
 │                                                                         │ 
 │   \ _ /    🔥 Temp:                      15.6°C                         │ 
 │  - ( )--o  💨 Wind:                      18.7 Km/h (250°)               │ 
 │   / (___)  🌈 Condition:                 Partly Cloudy                  │ 
 │            💧 Humidity:                  53.0%                          │ 
 │            🔥 Feels Like:                14.3°C                         │ 
 │            🧭 Surface Pressure:          1015.10 hPa                    │ 
 │            🌊 Sealevel Pressure:         1017.80 hPa                    │ 
 │            🌞 UV Index:                  6                              │ 
 │                                                                         │ 
 ╰────────────────────────────────────────────────────────────────────────╯ 
```

### 7-Day Forecast
```bash
./weather-Cli get London -f
```
```text
 ╭────────────────────────────────────────────────────────────────────────╮ 
 │                                                                         │ 
 │  WEATHER REPORT                                                         │ 
 │                                                                         │ 
 │ 📍 Location:                  London, United Kingdom                    │ 
 │ 🌐 Latitude:                  51.50853                                  │ 
 │ 🌐 Longitude:                 -0.12574                                  │ 
 │ 🕒 Timezone:                  Europe/London                             │ 
 │ 👥 Population:                8,961,989                                 │ 
 │                                                                         │ 
 │  CURRENT CONDITIONS                                                     │ 
 │                                                                         │ 
 │   \ _ /    🔥 Temp:                      15.6°C                         │ 
 │  - ( )--o  💨 Wind:                      18.7 Km/h (250°)               │ 
 │   / (___)  🌈 Condition:                 Partly Cloudy                  │ 
 │            💧 Humidity:                  53.0%                          │ 
 │            🔥 Feels Like:                14.3°C                         │ 
 │            🧭 Surface Pressure:          1015.10 hPa                    │ 
 │            🌊 Sealevel Pressure:         1017.80 hPa                    │ 
 │            🌞 UV Index:                  6                              │ 
 │                                                                         │ 
 │  7-DAY FORECAST                                                         │ 
 │                                                                         │ 
 │ 🗓️  Wed, May 20:              18.6/12.8°C     Light Drizzle             │ 
 │ 🗓️  Thu, May 21:              22.8/12.3°C     Overcast                  │ 
 │ 🗓️  Fri, May 22:              26.6/16.9°C     Overcast                  │ 
 │ 🗓️  Sat, May 23:              26.9/19.6°C     Overcast                  │ 
 │ 🗓️  Sun, May 24:              31.3/18.9°C     Clear Sky                 │ 
 │ 🗓️  Mon, May 25:              29.0/20.9°C     Clear Sky                 │ 
 │ 🗓️  Tue, May 26:              20.6/13.1°C     Clear Sky                 │ 
 │                                                                         │ 
 ╰────────────────────────────────────────────────────────────────────────╯
```

### Raw JSON Output (Automation-Friendly)
```bash
./weather-Cli get London --raw
```
```json
{
  "location": {
    "name": "London",
    "country": "United Kingdom",
    "latitude": 51.50853,
    "longitude": -0.12574,
    "timezone": "Europe/London",
    "population": 8961989
  },
  "forecast": {
    "current_weather": {
      "temperature": 15.6,
      "windspeed": 18.7,
      "winddirection": 250,
      "weathercode": 2
    },
    ...
  }
}
```

## Configuration ⚙️
You can manage your configuration using the `config` command:

```bash
# Set a default city
./weather-Cli config set default_city "Amsterdam"

# Disable interactive selection
./weather-Cli config set interactive false

# Set default units
./weather-Cli config set units imperial
```

The configuration is stored at `$HOME/.weather-cli.yaml`.

## Flags 🚩
| Flag | Shorthand | Description |
| :--- | :--- | :--- |
| `--forecast` | `-f` | Show 7-day extended forecast |
| `--hourly` | `-H` | Show 24-hour detailed forecast |
| `--units` | `-u` | Units to use (`metric` or `imperial`) |
| `--theme` | `-t` | Theme to use (`vibrant`, `dark`, `light`) |
| `--output` | `-o` | Export format (`json`, `csv`, `markdown`) |
| `--output-file` | | Path to save the exported data |
| `--no-list` | `-l` | Disable interactive city selection |
| `--no-cache` | | Disable local caching |
| `--no-style` | | Disable styled terminal visuals |
| `--raw` | `-r` | Output valid, unified JSON |
| `--config` | | Path to a custom config file |
| `--version` | `-v` | Print version number |

## Acknowledgments 📚
* API used: [Open-Meteo](https://open-meteo.com/)
* Built with: [Cobra](https://github.com/spf13/cobra), [Viper](https://github.com/spf13/viper), [Lipgloss](https://github.com/charmbracelet/lipgloss), and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Questions ⁉️
### If you have any questions or just want to talk, add me on [Discord](https://rayrsn.me/LinkInBio) or use my username **Rayr#6401**.
