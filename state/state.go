package state

type State string

const (
	Start              = "Start"
	CountWait          = "CountWait"
	FirstCurrencyWait  = "FirstCurrencyWait"
	SecondCurrencyWait = "SecondCurrencyWait"
)
