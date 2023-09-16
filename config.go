package main

const (
	StartText             = "_Ahoy there, bold soul \\! 'Tis I, Weather Report, your guide through the whims of meteorological fate\\. Ready yourself to brave the elements and seize control of your destiny\\! Share with me your coordinates, and I shall unravel the celestial tapestry, revealing the forecast  beyond even the grasp of Stands\\.🌦️🌪️_"
	Start                 = "/start"
	scheme                = "https"
	host                  = "api.openweathermap.org"
	pathWeather           = "data/2.5/weather"
	pathForecast          = "data/2.5/forecast"
	CurrentWeather        = "current weather"
	FiveDay               = "forecast for 5 days"
	ChoiceForecast        = "_Before we embark on this meteorological journey, dear traveler, tell me: do you seek the current weather, or does your curiosity extend to a forecast spanning five days? Your choice shall shape the course of our weather Stand's guidance\\!🌞_"
	ShareGeoloc           = "Share coordinates, traveler.🌐🌦️"
	FinishText            = "_Prepare for the heat and the chill, for the weather's Stand is as unpredictable as fate itself\\! Stay resolute, intrepid soul\\._⚔️"
	CurrentReportTemplate = "*__🌦️ Weather Stand Report 🌦️__*\n**Location:** %s📍\n" +
		"*Weather:* %s🌐\n*Temperature:* %.0f °C🌡\n️*Feels Like:* %.0f °C\n*Pressure:*%.0f mmHg\n*Humidity:* %d %%💧\n*Wind:* %.0f m/s💨%v\n" +
		"*Clouds:* %d %%☁️%v\n*Sunrise:* %v🌄\n*Sunset:* %s🌅\n*Reported at:* %v ⏰\n%s"
	RainTemplate = "\n*Rain:* %.2f mm/hr☔"
	SnowTemplate = "\n*Snow:* %.2f mm/hr☃️"
)

var latitude, longitude float64

type WeatherResponse struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Rain struct {
		OneH float64 `json:"1h"`
	} `json:"rain"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Snow struct {
		OneH   float64 `json:"1h"`
		ThreeH float64 `json:"3h"`
	} `json:"snow"`
	Dt  int `json:"dt"`
	Sys struct {
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	Name     string `json:"name"`
}
type ForecastResponse struct {
	List []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Gust  float64 `json:"gust"`
		} `json:"wind"`
		Visibility int    `json:"visibility"`
		Pop        int    `json:"pop"`
		DtTxt      string `json:"dt_txt"`
	} `json:"list"`
	City struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"city"`
}
