package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Error loading .env file")
		os.Exit(1)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false, TimeFormat: "2006-01-02 15:04:05"})

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic()
	}

	bot.Debug = true

	log.Info().Caller().Msgf("Authorized on account: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil { // ignore non-Message updates
			continue
		}
		log.Info().Caller().Msgf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, StartText)
		msg.ParseMode = "MarkdownV2"
		switch {
		case update.Message.Text == Start:
			msg.Text = StartText
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButtonLocation(ShareGeoloc),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(NeedHelp),
				),
			)
		case update.Message.Text == Help || update.Message.Text == NeedHelp:
			msg.Text = HelpMessage
		case update.Message.Location != nil:
			msg.Text = ChoiceForecast
			latitude = update.Message.Location.Latitude
			longitude = update.Message.Location.Longitude
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(CurrentWeather),
					tgbotapi.NewKeyboardButton(FiveDay),
					tgbotapi.NewKeyboardButtonLocation(ShareNewGeoloc),
				),
			)
		case update.Message.Text == CurrentWeather:
			msg.Text = RequestCurrentWeather(latitude, longitude)
		case update.Message.Text == FiveDay:
			msg.Text = RequestForecast(latitude, longitude)
		default:

		}
		msg.ReplyToMessageID = update.Message.MessageID
		if _, err := bot.Send(msg); err != nil {
			log.Error().Caller().Err(err).Msg("Error post message")
		}

	}

}
func RequestCurrentWeather(lat, lon float64) string {
	//https://api.openweathermap.org/data/2.5/weather?lat={lat}&lon={lon}&appid={API key}
	params := url.Values{
		"lat":   []string{strconv.FormatFloat(lat, 'f', -1, 64)},
		"lon":   []string{strconv.FormatFloat(lon, 'f', -1, 64)},
		"units": []string{"metric"},
		"appid": []string{os.Getenv("APIKEY")},
	}

	u := &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     pathWeather,
		RawQuery: params.Encode(),
	}
	resp, err := http.Get(u.String())
	if err != nil {
		log.Error().Caller().Err(err).Msg("Error get weather request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Caller().Err(err).Msg("Error reading resp.Body")
	}
	log.Debug().Msgf("%s", body)
	var Weather WeatherResponse
	err = json.Unmarshal(body, &Weather)
	if err != nil {
		log.Error().Caller().Err(err).Msg("Error unmarshalling json")
	}
	var rain, snow string
	if Weather.Rain.OneH != 0 || Weather.Snow.OneH != 0 {
		rain = strings.Replace(fmt.Sprintf(RainTemplate, Weather.Rain.OneH), ".", ",", -1)
		snow = strings.Replace(fmt.Sprintf(SnowTemplate, Weather.Snow.OneH), ".", ",", -1)
	}
	pressure := float64(Weather.Main.Pressure) / PascalsToHgmm
	WeatherMessage := fmt.Sprintf(CurrentReportTemplate, Weather.Name, Weather.Weather[0].Description, Weather.Main.Temp, Weather.Main.FeelsLike, pressure, Weather.Main.Humidity,
		Weather.Wind.Speed, rain, Weather.Clouds.All, snow, IntToTime(Weather.Sys.Sunrise), IntToTime(Weather.Sys.Sunset), IntToTime(Weather.Dt), FinishCurrentMessage)
	log.Info().Msgf("%s", WeatherMessage)
	return WeatherMessage
}
func RequestForecast(lat, lon float64) string {
	// api.openweathermap.org/data/2.5/forecast?lat={lat}&lon={lon}&appid={API key}
	params := url.Values{
		"lat":   []string{strconv.FormatFloat(lat, 'f', -1, 64)},
		"lon":   []string{strconv.FormatFloat(lon, 'f', -1, 64)},
		"units": []string{"metric"},
		"appid": []string{os.Getenv("APIKEY")},
	}

	u := &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     pathForecast,
		RawQuery: params.Encode(),
	}
	resp, err := http.Get(u.String())
	if err != nil {
		log.Error().Caller().Err(err).Msg("Error get weather request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Caller().Err(err).Msg("Error reading resp.Body")
	}
	log.Info().Msgf("%s", body)
	var Forecast ForecastResponse
	err = json.Unmarshal(body, &Forecast)
	if err != nil {
		log.Error().Caller().Err(err).Msg("Error unmarshalling json")
	}
	return CreateForecast(Forecast)
}
func IntToTime(n int) string {
	return time.Unix(int64(n), 0).Format("15:04:05")
}
func CreateForecast(Forecasts ForecastResponse) string {
	today := time.Unix(int64(Forecasts.List[0].Dt), 0)

	weather := []DayData{{
		Dt:       today.Format("Monday, 02-Jan-06"),
		Temp:     []float64{},
		Pressure: []float64{},
		Humidity: []float64{},
		Wind:     []float64{},
		Pop:      []float64{},
	},
	}
	j := 0
	for _, forecast := range Forecasts.List {
		date := time.Unix(int64(forecast.Dt), 0)
		if date.Day() != today.Day() {
			today = date
			j++
			weather = append(weather, struct {
				Dt       string
				Temp     []float64
				Pressure []float64
				Humidity []float64
				Wind     []float64
				Pop      []float64
			}{})
			weather[j].Dt = date.Format("Monday, 02-Jan-06")
		}
		weather[j].Temp = append(weather[j].Temp, forecast.Main.Temp)
		weather[j].Pressure = append(weather[j].Pressure, forecast.Main.Pressure)
		weather[j].Humidity = append(weather[j].Humidity, forecast.Main.Humidity)
		weather[j].Wind = append(weather[j].Wind, forecast.Wind.Speed)
		weather[j].Pop = append(weather[j].Pop, forecast.Pop)
	}
	fmt.Println(weather)
	var DaysWeather []DayWeather
	for _, day := range weather {
		DaysWeather = append(DaysWeather, struct {
			Dt       string
			TempMin  float64
			TempMax  float64
			Pressure float64
			Humidity float64
			Wind     float64
			Pop      float64
		}{Dt: day.Dt, TempMin: slices.Min(day.Temp), TempMax: slices.Max(day.Temp), Pressure: Average(day.Pressure) / PascalsToHgmm, Humidity: Average(day.Humidity), Wind: Average(day.Wind), Pop: Average(day.Pop)})
	}
	fmt.Println(DaysWeather)
	return CreateMessage(DaysWeather, Forecasts)
}
func CreateMessage(DaysWeather []DayWeather, Forecast ForecastResponse) string {
	Message := fmt.Sprintf(CityTemplate, Forecast.City.Name)
	for _, day := range DaysWeather {
		Message += fmt.Sprintf(DaysReportTemplate, day.Dt, day.TempMax, day.TempMin, day.Wind, day.Humidity, day.Pressure, day.Pop)
	}
	Message = strings.Replace(Message, "-", "\\-", -1)
	Message = strings.Replace(Message, ".", "\\.", -1)
	return Message + FinishDaysMessage
}
func Average(values []float64) float64 {
	total := 0.0
	for _, val := range values {
		total += val
	}
	return total / float64(len(values))
}
