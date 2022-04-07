package entities

import cs "xendit/shared/constants"

type Transaction struct {
	Id            int64               `json:"id"`
	UserID        string              `json:"user_id"`
	CardID        string              `json:"card_id"`
	TransactionID string              `json:"transaction_id"`
	Amount        int64               `json:"amount"`
	Created       int64               `json:"created"`
	LastModified  int64               `json:"last_modified"`
	State         cs.TransactionState `json:"state"`
}
