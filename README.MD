<div align="center">
  
# `🌦 Weather-Cli`
## **An app made to get weather information through the terminal**

<p align="center">
  <img width=40% height=40% src="https://github.com/Rayrsn/Weather-Cli/raw/main/cloud.gif?raw=true" alt="cloud">
</p>

[![Status Badge](https://github.com/Rayrsn/Weather-Cli/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/Rayrsn/Weather-Cli/actions/workflows/build.yml)

</div>

## Features 🏆
* Can show 
  * Population 
  * Latitude
  * Longitude
  * Weather
    * Temperature
    * Humidity
    * Wind Speed
    * Wind Direction
    * Weather Conditions
    * Feels Like Temperature (Apparent temperature based on a combination of factors e.g. humidity, wind speed, solar radiation, wind chill and air pressure)
    * Surface Pressure
    * Sealevel Pressure
    * UV Index
* Raw JSON output

## Building 🔨
* If you don't want to manually build the app, you can use:
  * The [weather-cli](https://aur.archlinux.org/packages/weather-cli) AUR if you're on Arch
  * Grab a binary from the [Artifacts](https://github.com/Rayrsn/Weather-Cli/actions/runs/2855792914)

Otherwise, you can build the program yourself by running:

```bash
go build
```
* **The executable will be located in the same directory you ran the command from**

## Usage 💻
* Example 1:
```bash
./weather-Cli get Tokyo
__________________________
Searching for city Tokyo...

City/Country: Tokyo/Japan
Latitude: 35.689500
Longitude: 139.691710
Timezone: Asia/Tokyo
Population: 8,336,599 (8.336599e+06)

Weather Info:
        Temperature: 28.5°
        Wind Direction: 168°
        Wind Speed: 7.0 Km/h
        Weather Condition: Partly Cloudy
        Humidity: 50.00%
        Real Feel: 39.1°
        Surface Pressure: 999.60 hPa
        Sealevel Pressure: 1004.00 hPa
        UV Index: 6
```
* Example 2:
```bash
./weather-Cli get Tokyo --raw
__________________________
(Returns the same data but in JSON)
```

## Acknowledgments and FAQ 📚
API used: [Open-Meteo](https://open-meteo.com/)

Built with: [Cobra](https://github.com/spf13/cobra)

## Questions ⁉️
### If yall have any questions or just wanna talk, add me on [Discord](https://rayr.ml/LinkInBio) or use my username Rayr#6401 (this might change so it's better to just use the link)
