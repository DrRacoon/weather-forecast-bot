package main

const (
	PascalsToHgmm         = 1.333
	StartText             = "_Ahoy there, bold soul \\! 'Tis I, Weather Report, your guide through the whims of meteorological fate\\. Ready yourself to brave the elements and seize control of your destiny\\! Share with me your coordinates, and I shall unravel the celestial tapestry, revealing the forecast  beyond even the grasp of Stands\\.🌦️🌪️_"
	Start                 = "/start"
	Help                  = "/help"
	NeedHelp              = "Yare Yare Daze, I need help!💪🌟"
	HelpMessage           = "_Make sure your phone's geolocation is enabled, and the necessary settings for Telegram are configured\\. Alternatively, you can manually set your location on the map \\(click on the paperclip icon, then choose 'Location'\\)\\. Stand strong, and let us embark on this weather journey\\!🌍📍_"
	scheme                = "https"
	host                  = "api.openweathermap.org"
	pathWeather           = "data/2.5/weather"
	pathForecast          = "data/2.5/forecast"
	CurrentWeather        = "current weather"
	FiveDay               = "forecast for 5 days"
	ChoiceForecast        = "_Before we embark on this meteorological journey, dear traveler, tell me: do you seek the current weather, or does your curiosity extend to a forecast spanning five days? Your choice shall shape the course of our weather Stand's guidance\\!🌞_"
	ShareGeoloc           = "Share coordinates, traveler.🌐🌦️"
	ShareNewGeoloc        = "share new coordinates"
	FinishCurrentMessage  = "_Prepare for the heat and the chill, for the weather's Stand is as unpredictable as fate itself\\! Stay resolute, intrepid soul\\._⚔️"
	CurrentReportTemplate = "*__🌦️ Weather Stand Report 🌦️__*\n*Location:* %23s📍\n" +
		"*Weather:* %23s 🌐\n*Temperature:* %15.0f°C 🌡\n️*Feels Like:* %20.0f°C\n*Pressure:*%19.0f mmHg\n*Humidity:* %21d%% 💧\n*Wind:*%25.0f m/s 💨%v\n" +
		"*Clouds:* %24d%% ☁️%v\n*Sunrise: *%22v 🌄\n*Sunset:* %23s 🌅\n*Reported at:* %14v ⏰\n%s"
	RainTemplate       = "\n*Rain:* %18.2f mm/hr ☔"
	SnowTemplate       = "\n*Snow:* %16.2f mm/hr ☃️"
	CityTemplate       = "_Behold, the climatic predictions for this splendid __%s__\\!_🌦\n"
	DaysReportTemplate = "__%s__\n*Temperature:*\n\t*• Day:*%.1f°C 🌞\n\t*• Night:* %.1f°C 🌛\n*Wind:* %.1f m/s 💨\n*Humidity:* %.0f%% 💧\n*Pressure:* %.0f mmHg 🌡️\n*Precipitation probability:* %.1f%% 🌦️❄️\n\n"
	FinishDaysMessage  = "_Prepare for the weather adventure, dear traveler\\! May the winds be in your favor and the skies clear on your path\\. Until we meet again, stand resilient and embrace each day's forecast with the heart of a true weather warrior\\!_ ⚔️🌦️🌪️"
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
			Pressure  float64 `json:"pressure"`
			Humidity  float64 `json:"humidity"`
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
		Pop   float64 `json:"pop"`
		DtTxt string  `json:"dt_txt"`
	} `json:"list"`
	City struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"city"`
}
type DayWeather struct {
	Dt       string
	TempMin  float64
	TempMax  float64
	Pressure float64
	Humidity float64
	Wind     float64
	Pop      float64
}
type DayData struct {
	Dt       string
	Temp     []float64
	Pressure []float64
	Humidity []float64
	Wind     []float64
	Pop      []float64
}
