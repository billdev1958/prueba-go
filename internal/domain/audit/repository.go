package audit

import "context"

type AuditRepository interface {
	Save(ctx context.Context, log *AuditLog) error
	GetAll(ctx context.Context) ([]AuditLog, error)
}
