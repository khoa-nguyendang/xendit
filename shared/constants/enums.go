package constants

type CardState int
type UserState int
type TransactionState int

const (
	USER_STATE_DEFAULT UserState = iota
	USER_STATE_BLOCKED
	USER_STATE_LIMITED
)

const (
	CARD_STATE_DEFAULT CardState = iota
	CARD_STATE_BLOCKED
	CARD_STATE_LIMITED
)

const (
	TRANSACTION_STATE_DEFAULT TransactionState = iota
	TRANSACTION_STATE_SUCCESSED
	TRANSACTION_STATE_FAILED
)
