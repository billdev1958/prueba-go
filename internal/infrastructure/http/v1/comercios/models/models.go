package models

import "time"

type CreateComercioRequest struct {
	Name          string `json:"name" binding:"required" example:"Abarrotes Doña María"`
	ComissionRate string `json:"comission_rate" binding:"required" example:"0.025"`
}

type UpdateComercioRequest struct {
	Name          string `json:"name" example:"Abarrotes Doña María Sucursal 2"`
	ComissionRate string `json:"comission_rate" example:"0.030"`
}

type ComercioResponse struct {
	ID            string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name          string    `json:"name" example:"Abarrotes Doña María"`
	ComissionRate string    `json:"comission_rate" example:"0.025"`
	CreatedAt     time.Time `json:"created_at" example:"2023-10-27T10:00:00Z"`
}
