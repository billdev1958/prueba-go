package models

import (
	"prueba-go/pkg/types"
	"time"
)

type AuditResponse struct {
	LogID      types.UID `json:"log_id" example:"a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d"`
	Action     string    `json:"action" example:"COMMERCE_CREATE"`
	Actor      string    `json:"actor" example:"admin_user"`
	ResourceID types.UID `json:"resource_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Timestamp  time.Time `json:"timestamp" example:"2023-10-27T10:00:00Z"`
}
