<div align="center">
  
# `🌦 Weather-Cli`
## **An app made to get weather information through the terminal**

<p align="center">
  <img width=40% height=40% src="https://github.com/Rayrsn/Weather-Cli/raw/main/cloud.gif?raw=true" alt="MAIN_IMAGE_NAME">
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
* Raw JSON output


## Building 🔨
```bash
go build
```
* **The executable will be located in the same directory you ran the command from**

## Usage 💻
* Example 1:
```bash
./weather-cli Tokyo
__________________________
Searching for city Tokyo ...

City/Country: Tokyo/Japan
Latitude: 35.689500
Longitude: 139.691710
Timezone: Asia/Tokyo
Population: 8,336,599 (8.336599e+06)

Weather Info:
        Temperature: 27.4°
        Wind Direction: 192°
        Wind Speed: 15.5 Km/h
        Weather Condition: Overcast
        Humidity: 85.91%

```
* Example 2:
```bash
./weather-cli Tokyo --raw
__________________________
(Prints the same data but in JSON)
```

## Acknowledgments and FAQ 📚
Api used : [Open-Meteo](https://open-meteo.com/)

## Questions ⁉️
### If yall have any questions or just wanna talk, add me on [Discord](https://rayr.ml/LinkInBio) or use my username Rayr#6401 (this might change so it's better to just use the link)
