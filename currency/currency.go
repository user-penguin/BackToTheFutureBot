package currency

type Currency string

const (
	Euro      = "EUR"
	DollarUSA = "USD"
	Ruble     = "RUB"
)

func GetAllCurrencies() [3]string {
	return [...]string{
		Euro,
		DollarUSA,
		Ruble,
	}
}
