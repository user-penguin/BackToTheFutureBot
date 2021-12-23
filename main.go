package main

import (
	"BackToTheFutureBot/command"
	"BackToTheFutureBot/currency"
	"BackToTheFutureBot/message"
	"BackToTheFutureBot/reader"
	"BackToTheFutureBot/scene"
	"BackToTheFutureBot/state"
	exchange "github.com/3crabs/go-yahoo-finance-api"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
)

type UserCondition struct {
	Scene        string
	State        string
	CurrencyFrom string
	Value        float64
	CurrencyTo   string
}

// мапа состояний пользователей
var states map[int64]UserCondition

type Yahoo struct {
	Token string
}

func (y *Yahoo) getRatio(s string, s2 string) (float64, error) {
	pair := exchange.Pair{From: s, To: s2}
	quote, err := exchange.GetCurrency(pair, y.Token)
	if err != nil {
		return 0, err
	}
	return quote.QuoteResponse.Result[0].Ask, nil
}

// в дальнейшем допилить работу с глобальной мапой
func getState(chatId int64) (UserCondition, bool) {
	userState, ok := states[chatId]
	return userState, ok
}

// в дальнейшем допилить работу с глобальной мапой
func setState(chatId int64, currency UserCondition) {
	states[chatId] = currency
}

func calculateSummary(count float64, from string, to string, yahoo *Yahoo) float64 {
	ratio, _ := yahoo.getRatio(from, to)
	return ratio * count
}

func main() {
	run()
}

func baseSceneHandle(chatId int64, text string, bot *tgbot.BotAPI) {
	switch text {
	case string(command.Start):
		setState(chatId, UserCondition{State: state.Begin, Scene: scene.MainMenu})
		msg := tgbot.NewMessage(chatId, message.OnStartButtonMessage())
		msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
		_, _ = bot.Send(msg)
	default:
		msg := tgbot.NewMessage(chatId, message.DefaultStartButtonMessage())
		msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
		_, _ = bot.Send(msg)
	}
}

func mainMenuSceneHandle(chatId int64, text string, bot *tgbot.BotAPI) {
	switch text {
	case string(command.Convert):
		setState(chatId, UserCondition{State: state.FirstCurrencyWait, Scene: scene.ConvertMenu})
		msg := tgbot.NewMessage(chatId, message.OnConvertButtonMessage())
		msg.ReplyMarkup = message.GetCurrenciesKeyboard()
		_, _ = bot.Send(msg)
	case string(command.Start):
		msg := tgbot.NewMessage(chatId, message.OnStartButtonMessage())
		msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
		_, _ = bot.Send(msg)
	default:
		msg := tgbot.NewMessage(chatId, message.WrongCommandMessage())
		msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
		_, _ = bot.Send(msg)
	}
}

func convertSceneHandle(chatId int64, text string, bot *tgbot.BotAPI, yahoo *Yahoo) {
	userCondition, ok := getState(chatId)
	if !ok {
		states[chatId] = UserCondition{State: state.Begin, Scene: "/start"}
		userCondition = states[chatId]
	}
	switch userCondition.State {
	// text здесь - это наименование валюты
	case state.FirstCurrencyWait:
		_, ok := currency.GetCurrencyMap()[text]
		if !ok {
			msg := tgbot.NewMessage(chatId, message.WrongCurrencyMessage())
			_, _ = bot.Send(msg)
		} else {
			userCondition.State = state.CountWait
			userCondition.CurrencyFrom = text
			states[chatId] = userCondition
			msg := tgbot.NewMessage(chatId, message.TypeCountOfCurrencyMessage())
			msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
			_, _ = bot.Send(msg)
		}
	case state.CountWait:
		countMoney, err := strconv.ParseFloat(text, 64)
		if err != nil {
			log.Println(err)
			msg := tgbot.NewMessage(chatId, message.WrongCountOfCurrencyMessage())
			msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
			_, _ = bot.Send(msg)
		} else {
			userCondition.State = state.SecondCurrencyWait
			userCondition.Value = countMoney
			states[chatId] = userCondition
			msg := tgbot.NewMessage(chatId, message.SelectSecondCurrencyMessage())
			msg.ReplyMarkup = message.GetCurrenciesKeyboard()
			_, _ = bot.Send(msg)
		}
	case state.SecondCurrencyWait:
		_, ok := currency.GetCurrencyMap()[text]
		if !ok {
			msg := tgbot.NewMessage(chatId, message.WrongCurrencyMessage())
			_, _ = bot.Send(msg)
		} else {
			userCondition.State = state.Begin
			userCondition.Scene = scene.MainMenu
			userCondition.CurrencyTo = text
			states[chatId] = userCondition
			resultValue := calculateSummary(userCondition.Value, userCondition.CurrencyFrom, userCondition.CurrencyTo, yahoo)
			msg := tgbot.NewMessage(chatId, message.ResultConvertMessage(userCondition.CurrencyFrom, userCondition.CurrencyTo, userCondition.Value, resultValue))
			msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
			_, _ = bot.Send(msg)
		}
	}
}

func run() {
	config, err := reader.GetConfig()
	if err != nil {
		panic(err)
	}
	bot, err := tgbot.NewBotAPI(config.BotToken)
	if err != nil {
		log.Println(err)
		return
	}
	yahoo := &Yahoo{Token: config.YahooToken}

	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	states = make(map[int64]UserCondition)

	log.Println("Bot is start up!")

	for update := range updates {
		// empty message
		if update.Message == nil {
			continue
		}

		text := update.Message.Text
		chatId := update.Message.Chat.ID

		userCondition, ok := getState(chatId)
		if !ok {
			states[chatId] = UserCondition{State: state.Begin, Scene: scene.Base}
			userCondition = states[chatId]
		}

		switch userCondition.Scene {
		case scene.Base:
			baseSceneHandle(chatId, text, bot)
		case scene.MainMenu:
			mainMenuSceneHandle(chatId, text, bot)
		case scene.ConvertMenu:
			convertSceneHandle(chatId, text, bot, yahoo)
		}
	}

}
