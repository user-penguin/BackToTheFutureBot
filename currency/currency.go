package currency

type Currency string

const (
	Euro      = "EUR"
	DollarUSA = "USD"
	Ruble     = "RUB"
)

func GetCurrencyMap() map[string]interface{} {
	curMap := make(map[string]interface{}, len(GetAllCurrencies()))
	for _, currency := range GetAllCurrencies() {
		curMap[currency] = true
	}
	return curMap
}

func GetAllCurrencies() [3]string {
	return [...]string{
		Euro,
		DollarUSA,
		Ruble,
	}
}
