package model

type NewTransactionResponse struct {
	Message string `json:"message"`
}

type TransactionRequest struct {
	AddressRecipient string  `json:"addressRecipient"`
	Sum              float64 `json:"sum"`
}
