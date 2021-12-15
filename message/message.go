package message

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
