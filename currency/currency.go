package currency

import exchange "github.com/3crabs/go-yahoo-finance-api"

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

func GetAllPair() []exchange.Pair {
	pairs := [...]exchange.Pair{
		{Ruble, Euro},
		{Euro, Ruble},
		{Ruble, DollarUSA},
		{DollarUSA, Ruble},
		{DollarUSA, Euro},
		{Euro, Ruble},
	}
	return pairs[:]
}
