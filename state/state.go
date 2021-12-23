package state

type State string

// условимся, что состояние обозначает момент сейчас
// то есть юзер сейчас стоит перед нажатием на /конверт
// он в состоянии бегин, в дефолтном состоянии
// - пользователь нажал /конверт, его сразу перебросило в FirstCurrencyWait
//   тут же ему выбрасывается сообщения с кнопками валют
const (
	Begin              = "begin" // default value means the User in main menu
	CountWait          = "count-wait"
	FirstCurrencyWait  = "first-currency-wait"
	SecondCurrencyWait = "second-currency-wait"
	End                = "end"
)
