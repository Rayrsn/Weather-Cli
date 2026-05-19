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
