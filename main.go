package main

import (
	"BackToTheFutureBot/command"
	"BackToTheFutureBot/currency"
	"BackToTheFutureBot/message"
	"BackToTheFutureBot/reader"
	"BackToTheFutureBot/state"
	"errors"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
)

type Currency struct {
	CurrencyFrom string
	CurrencyTo   string
	Value        int
	State        string
}

var states map[int64]Currency

// в дальнейшем допилить работу с глобальной мапой
func getState(chatId int64) (Currency, bool) {
	currency, ok := states[chatId]
	return currency, ok
}

// в дальнейшем допилить работу с глобальной мапой
func setState(chatId int64, currency Currency) {
	states[chatId] = currency
}

func getMsgByState(chatId int64, message string) (string, error) {
	currency, ok := getState(chatId)
	if !ok {
		states[chatId] = Currency{State: "start"}
		currency = states[chatId]
	}

	switch currency.State {
	case state.Start:
		currency.State = state.FirstCurrencyWait
		states[chatId] = currency
		return "Привет. Из какой валюты будем конвертировать?", nil
	case state.FirstCurrencyWait:
		currency.State = state.CountWait
		currency.CurrencyFrom = message
		states[chatId] = currency
		return "Сколько:", nil
	case state.CountWait:
		currency.State = state.SecondCurrencyWait
		currency.Value, _ = strconv.Atoi(message)
		states[chatId] = currency
		return "Куда?", nil
	case state.SecondCurrencyWait:
		currency.State = state.Start
		currency.CurrencyTo = message
		states[chatId] = currency
		return "вы получите кучу денег!" + strconv.Itoa(currency.Value*20), nil
	default:
		return "", errors.New("ввели какую-то херню")
	}
}

func main() {
	run()
}

func run() {
	token, err := reader.GetTokenFromFile()
	if err != nil {
		panic(err)
	}
	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		log.Println(err)
		return
	}
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
		currency, ok := getState(chatId)
		if !ok {
			states[chatId] = Currency{State: state.Start}
			currency = states[chatId]
		}

		switch text {
		case string(command.Start):
			setState(chatId, Currency{State: state.Start})
			msg := tgbot.NewMessage(chatId, message.StartMessage())
			msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
			_, _ = bot.Send(msg)
		case string(command.Convert):
			setState(chatId, Currency{State: state.FirstCurrencyWait})
			msg := tgbot.NewMessage(chatId, message.SelectFirstCurrency())
			msg.ReplyMarkup = numericKeyboard
			_, _ = bot.Send(msg)
		default:
			textMessage, err := getMsgByState(chatId, text)
			if err != nil {
				log.Println("Они ввели какую-то херню:", text, chatId, " state: ", currency.State)
				continue
			}
			msg := tgbot.NewMessage(chatId, textMessage)
			if currency.State == state.SecondCurrencyWait || currency.State == state.FirstCurrencyWait {
				msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
			} else {
				msg.ReplyMarkup = numericKeyboard
			}
			_, _ = bot.Send(msg)
		}
	}

}
