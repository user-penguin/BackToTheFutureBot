package main

import (
	"BackToTheFutureBot/command"
	"BackToTheFutureBot/currency"
	"BackToTheFutureBot/message"
	"BackToTheFutureBot/reader"
	"BackToTheFutureBot/state"
	"errors"
	"fmt"
	exchange "github.com/3crabs/go-yahoo-finance-api"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
)

type Currency struct {
	CurrencyFrom string
	CurrencyTo   string
	Value        float64
	State        string
}

type Yahoo struct {
	Token string
}

func (y *Yahoo) getRatio(s string, s2 string) (float64, error) {
	quote, err := exchange.GetCurrency(s, s2, y.Token)
	if err != nil {
		return 0, err
	}
	return quote.QuoteResponse.Result[0].Ask, nil
}

var states map[int64]Currency

// в дальнейшем допилить работу с глобальной мапой
func getState(chatId int64) (Currency, bool) {
	userState, ok := states[chatId]
	return userState, ok
}

// в дальнейшем допилить работу с глобальной мапой
func setState(chatId int64, currency Currency) {
	states[chatId] = currency
}

func getMsgByState(chatId int64, message string, yahoo *Yahoo) (string, error) {
	userCondition, ok := getState(chatId)
	if !ok {
		states[chatId] = Currency{State: "start"}
		userCondition = states[chatId]
	}

	switch userCondition.State {
	case state.Begin:
		userCondition.State = state.FirstCurrencyWait
		states[chatId] = userCondition
		return "Привет. Из какой валюты будем конвертировать?", nil
	case state.FirstCurrencyWait:
		userCondition.State = state.CountWait
		userCondition.CurrencyFrom = message
		states[chatId] = userCondition
		return "Сколько:", nil
	case state.CountWait:
		userCondition.State = state.SecondCurrencyWait
		userCondition.Value, _ = strconv.ParseFloat(message, 64)
		states[chatId] = userCondition
		return "Куда?", nil
	case state.SecondCurrencyWait:
		userCondition.State = state.Begin
		userCondition.CurrencyTo = message
		states[chatId] = userCondition
		summary := calculateSummary(userCondition.Value, userCondition.CurrencyFrom, userCondition.CurrencyTo, yahoo)
		return fmt.Sprintf("Вы получите: %f %s", summary, userCondition.CurrencyTo), nil
	default:
		return "", errors.New("ввели какую-то херню")
	}
}

func calculateSummary(count float64, from string, to string, yahoo *Yahoo) float64 {
	ratio, _ := yahoo.getRatio(from, to)
	return ratio * count
}

func main() {
	run()
}

func standardMenuHandle(chatId int64, text string, bot *tgbot.BotAPI) {
	switch text {
	case string(command.Start):
		msg := tgbot.NewMessage(chatId, message.StartMessage())
		msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
		_, _ = bot.Send(msg)
	case string(command.Convert):
		setState(chatId, Currency{State: state.FirstCurrencyWait})
		msg := tgbot.NewMessage(chatId, message.SelectFirstCurrency())
		msg.ReplyMarkup = getCurrenciesKeyboard()
		_, _ = bot.Send(msg)
	}
}

func getCurrenciesKeyboard() tgbot.ReplyKeyboardMarkup {
	var keyboard [][]tgbot.KeyboardButton
	var row []tgbot.KeyboardButton
	for i, currencyName := range currency.GetAllCurrencies() {
		row = append(row, tgbot.NewKeyboardButton(currencyName))
		if (i+1)%3 == 0 {
			keyboard = append(keyboard, row)
			row = nil
		}
	}
	return tgbot.NewReplyKeyboard(keyboard...)
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
	states = make(map[int64]Currency)

	log.Println("Bot is start up!")

	var numericKeyboard = tgbot.NewReplyKeyboard(
		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton(currency.Euro),
			tgbot.NewKeyboardButton(currency.DollarUSA),
			tgbot.NewKeyboardButton(currency.Ruble),
		),
	)

	for update := range updates {
		// empty message
		if update.Message == nil {
			continue
		}
		text := update.Message.Text
		chatId := update.Message.Chat.ID

		// имеем кого-то и текст от этого кого-то
		// сначала нужно проверить состояние

		userCondition, ok := getState(chatId)
		if !ok {
			states[chatId] = Currency{State: state.Begin}
			userCondition = states[chatId]
		}

		if userCondition.State == state.Begin {
			standardMenuHandle(chatId, text, bot)
			continue
		}

		switch text {
		case string(command.Start):
			setState(chatId, Currency{State: state.Begin})
			msg := tgbot.NewMessage(chatId, message.StartMessage())
			msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
			_, _ = bot.Send(msg)
		case string(command.Convert):
			setState(chatId, Currency{State: state.FirstCurrencyWait})
			msg := tgbot.NewMessage(chatId, message.SelectFirstCurrency())
			msg.ReplyMarkup = numericKeyboard
			_, _ = bot.Send(msg)
		default:
			textMessage, err := getMsgByState(chatId, text, yahoo)
			if err != nil {
				log.Println("Они ввели какую-то херню:", text, chatId, " state: ", userCondition.State)
				continue
			}
			msg := tgbot.NewMessage(chatId, textMessage)
			if userCondition.State == state.SecondCurrencyWait || userCondition.State == state.FirstCurrencyWait {
				msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
			} else {
				msg.ReplyMarkup = numericKeyboard
			}
			_, _ = bot.Send(msg)
		}
	}

}
