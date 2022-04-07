package viewmodels

type TransactionResponse struct {
	Id                string      `json:"id"`
	RequestData       interface{} `json:"request_data"`
	Recommendation    string      `json:"recommendation"`
	Context           []string    `json:"context"`
	Created           string      `json:"created"`
	Updated           string      `json:"updated"`
	TransactionStatus string      `json:"transaction_status"`
}
