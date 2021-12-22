package message

import (
	"BackToTheFutureBot/currency"
	"fmt"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

func OnStartButtonMessage() string {
	msg := "Добро пожаловать в гости к боту-конвертеру!\n" +
		"Здесь вы можете:\n" +
		"- конвертировать валюту по команде [ /convert ]\n"
	return msg
}

func DefaultStartButtonMessage() string {
	return "Для начала работы введите команду [ /start ]"
}

func WrongCommandMessage() string {
	return "Ваша команда не распознана. Попробуйте ещё раз"
}

func WrongCurrencyMessage() string {
	return "К такой валюте мы ещё не готовы, следите за CHANGELOG в репозитории бота, попробуйте ещё раз"
}

func WrongCountOfCurrencyMessage() string {
	return "Вы ввели некорректное значение, попробуйте ещё раз"
}
func OnConvertButtonMessage() string {
	return "Отлично. Выберите, из какой валюты мы будем конвертировать:"
}

func TypeCountOfCurrencyMessage() string {
	return "Сколько денег вы хотите конвертировать?"
}

func SelectSecondCurrencyMessage() string {
	return "Выберите вторую валюту"
}

func ResultConvertMessage(curFrom string, curTo string, valFrom float64, valTo float64) string {
	return fmt.Sprintf("%f%s = %f%s", valFrom, curFrom, valTo, curTo)
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
