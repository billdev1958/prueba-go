package models

import (
	"prueba-go/pkg/types"
	"time"
)

type AuditResponse struct {
	LogID      types.UID `json:"log_id"`
	Action     string    `json:"action"`
	Actor      string    `json:"actor"`
	ResourceID types.UID `json:"resource_id"`
	Timestamp  time.Time `json:"timestamp"`
}
