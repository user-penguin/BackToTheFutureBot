package message

import (
	"BackToTheFutureBot/currency"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

func StartMessage() string {
	return "Добро пожаловать к боту-денежному конвертеру"
}

func ConvertMessage() string {
	return "Convert"
}

func SelectFirstCurrency() string {
	return "Окей, погнали. Выберите первую валюту"
}

func SelectSecondCurrency() string {
	return "Выберите вторую валюту"
}

func TypeCount(currency1 string, currency2 string) string {
	return "Сколько " + currency1 + " вы хотите конвертировать в " + currency2 + "?"
}

func GetCurrenciesKeyboard() tgbot.ReplyKeyboardMarkup {
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
