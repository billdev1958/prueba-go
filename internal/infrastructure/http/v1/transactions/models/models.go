package models

import "time"

type CreateTransactionRequest struct {
	CommercioID string `json:"comercio_id" binding:"required"`
	Amount      string `json:"amount" binding:"required"`
}

type TransactionResponse struct {
	ID          string    `json:"id"`
	CommercioID string    `json:"comercio_id"`
	Amount      string    `json:"amount"`
	AppliedRate string    `json:"commission"`
	Commission  string    `json:"fee"`
	NetAmount   string    `json:"net_amount"`
	CreatedAt   time.Time `json:"created_at"`
}
