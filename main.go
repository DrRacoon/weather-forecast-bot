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
			)
		case update.Message.Location != nil:
			msg.Text = ChoiceForecast
			latitude = update.Message.Location.Latitude
			longitude = update.Message.Location.Longitude
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(CurrentWeather),
					tgbotapi.NewKeyboardButton(FiveDay),
				),
			)
		case update.Message.Text == CurrentWeather:
			msg.Text = RequestCurrentWeather(latitude, longitude)
			//msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case update.Message.Text == FiveDay:
			msg.Text = RequestForecast(latitude, longitude)
			//msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
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
	log.Info().Msgf("%s", body)
	var Weather WeatherResponse
	err = json.Unmarshal(body, &Weather)
	if err != nil {
		log.Error().Caller().Err(err).Msg("Error unmarshalling json")
	}
	var rain, snow string
	if Weather.Rain.OneH != 0 {
		rain = strings.Replace(fmt.Sprintf(RainTemplate, Weather.Rain.OneH), ".", ",", -1)
	}
	if Weather.Snow.OneH != 0 {
		snow = strings.Replace(fmt.Sprintf(SnowTemplate, Weather.Snow.OneH), ".", ",", -1)

	}
	pressure := float64(Weather.Main.Pressure) / 1.333
	WeatherMessage := fmt.Sprintf(CurrentReportTemplate, Weather.Name, Weather.Weather[0].Description, Weather.Main.Temp, Weather.Main.FeelsLike, pressure, Weather.Main.Humidity,
		Weather.Wind.Speed, rain, Weather.Clouds.All, snow, IntToTime(Weather.Sys.Sunrise), IntToTime(Weather.Sys.Sunset), IntToTime(Weather.Dt), FinishText)
	fmt.Println(WeatherMessage)
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
	fmt.Println(Forecast)
	return GroupByForecast(Forecast)
}
func IntToTime(n int) string {
	return time.Unix(int64(n), 0).Format("15:04:05")
}
func GroupByForecast(Forecast ForecastResponse) string {

	return "we are working fo this"
}
