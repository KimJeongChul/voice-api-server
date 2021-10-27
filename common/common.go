package common

const (
	UNDEFINED = "UNDEFINED"
)

type SessionObj struct {
	FuncName      string
	TransactionId string
	Flush         func()
}
