package state

type State string

const (
	Begin              = "Begin" // default value means the User in main menu
	CountWait          = "CountWait"
	FirstCurrencyWait  = "FirstCurrencyWait"
	SecondCurrencyWait = "SecondCurrencyWait"
)
