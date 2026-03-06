package models

import "time"

type CreateTransactionRequest struct {
	CommercioID string `json:"comercio_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Amount      string `json:"amount" binding:"required" example:"100.00"`
}

type TransactionResponse struct {
	ID          string    `json:"id" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
	CommercioID string    `json:"comercio_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Amount      string    `json:"amount" example:"100.00"`
	AppliedRate string    `json:"commission" example:"0.025"`
	Commission  string    `json:"fee" example:"2.50"`
	NetAmount   string    `json:"net_amount" example:"97.50"`
	CreatedAt   time.Time `json:"created_at" example:"2023-10-27T10:05:00Z"`
}
