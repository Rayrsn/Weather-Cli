package api

type GeocodingResult struct {
	Name       string  `json:"name"`
	Country    string  `json:"country"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Timezone   string  `json:"timezone"`
	Population float64 `json:"population"`
}

type GeocodingResponse struct {
	Results []GeocodingResult `json:"results"`
}

type ForecastResponse struct {
	CurrentWeather struct {
		Temperature   float64 `json:"temperature"`
		Windspeed     float64 `json:"windspeed"`
		Winddirection float64 `json:"winddirection"`
		Weathercode   float64 `json:"weathercode"`
	} `json:"current_weather"`
	Hourly struct {
		Relativehumidity2M  []float64 `json:"relativehumidity_2m"`
		ApparentTemperature []float64 `json:"apparent_temperature"`
		SurfacePressure     []float64 `json:"surface_pressure"`
		PressureMsl         []float64 `json:"pressure_msl"`
	} `json:"hourly"`
	Daily struct {
		Time               []string  `json:"time"`
		Temperature2MMax   []float64 `json:"temperature_2m_max"`
		Temperature2MMin   []float64 `json:"temperature_2m_min"`
		Weathercode        []float64 `json:"weathercode"`
		Sunrise            []string  `json:"sunrise"`
		Sunset             []string  `json:"sunset"`
		UvIndexMax         []float64 `json:"uv_index_max"`
		PrecipitationSum   []float64 `json:"precipitation_sum"`
		Windspeed10MMax    []float64 `json:"windspeed_10m_max"`
		Winddirection10MDo []float64 `json:"winddirection_10m_dominant"`
	} `json:"daily"`
}

type AirQualityResponse struct {
	Hourly struct {
		UvIndex []float64 `json:"uv_index"`
	} `json:"hourly"`
}

type CombinedResponse struct {
	Location   GeocodingResult    `json:"location"`
	Forecast   ForecastResponse   `json:"forecast"`
	AirQuality AirQualityResponse `json:"air_quality"`
}
