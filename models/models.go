package models

type BalanceResponse struct {
	Balance string `json:"balance"`
}

type TransferRequest struct {
	To     string `json:"to"`
	Amount string `json:"amount"`
}