package viewmodels

type TransactionRequest struct {
	UserID        string  `json:"user_id"`
	CardID        string  `json:"card_id"`
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
}
