package audit

import (
	"prueba-go/pkg/types"
	"time"
)

type AuditLog struct {
	LogID      types.UID
	Action     string
	Actor      string
	ResourceID types.UID
	Timestamp  time.Time
}
