package models

import "time"

type CreateComercioRequest struct {
	Name          string `json:"name" binding:"required"`
	ComissionRate string `json:"comission_rate" binding:"required"`
}

type UpdateComercioRequest struct {
	Name          string `json:"name"`
	ComissionRate string `json:"comission_rate"`
}

type ComercioResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	ComissionRate string    `json:"comission_rate"`
	CreatedAt     time.Time `json:"created_at"`
}
