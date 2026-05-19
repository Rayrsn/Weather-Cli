<div align="center">
  
# `🌦 Weather-Cli`
## **A modern, fast CLI tool for weather information**

<p align="center">
  <img width=40% height=40% src="https://github.com/Rayrsn/Weather-Cli/raw/main/cloud.gif?raw=true" alt="cloud">
</p>

[![Status Badge](https://github.com/Rayrsn/Weather-Cli/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/Rayrsn/Weather-Cli/actions/workflows/build.yml)

</div>

## Features 🏆
* **Modern UI:** Beautiful, styled terminal output with emojis and colors using `lipgloss`.
* **Extended Forecast:** View a detailed 7-day outlook with the `--forecast` flag.
* **Global Support:** Geocoding for any city worldwide with full URL encoding support.
* **Unit Toggles:** Switch between `metric` (Celsius, Km/h) and `imperial` (Fahrenheit, mph).
* **High Performance:** Parallel API calls using Go routines for near-instant results.
* **Configuration:** Set a default city and preferred units in a config file.
* **Automation-Friendly:** Unified `--raw` JSON output for piping into tools like `jq`.
* **Reliable:** Struct-based API parsing and robust error handling.

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

### Get 7-Day Forecast
```bash
./weather-Cli get Paris --forecast
```

### Use Imperial Units
```bash
./weather-Cli get "New York" --units imperial
```

### Raw JSON Output
```bash
./weather-Cli get Tokyo --raw | jq .
```

### Disable Styling (Plain Text)
```bash
./weather-Cli get London --no-style
```

## Configuration ⚙️
You can create a configuration file at `$HOME/.weather-cli.yaml` to set your preferences:

```yaml
default_city: London
units: metric
no_style: false
```

## Flags 🚩
| Flag | Shorthand | Description |
| :--- | :--- | :--- |
| `--forecast` | `-f` | Show 7-day extended forecast |
| `--units` | `-u` | Units to use (`metric` or `imperial`) |
| `--raw` | `-r` | Output valid, unified JSON |
| `--no-style` | | Disable styled terminal visuals |
| `--config` | | Path to a custom config file |
| `--version` | `-v` | Print version number |

## Acknowledgments 📚
* API used: [Open-Meteo](https://open-meteo.com/)
* Built with: [Cobra](https://github.com/spf13/cobra), [Viper](https://github.com/spf13/viper), and [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Questions ⁉️
### If you have any questions or just want to talk, add me on [Discord](https://rayrsn.me/LinkInBio) or use my username **Rayr#6401**.
