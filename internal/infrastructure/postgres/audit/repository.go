package audit

import (
	"context"
	"prueba-go/internal/domain/audit"
	"prueba-go/internal/infrastructure/postgres"
)

type pgxAuditRepository struct {
	*postgres.PgxRepository
}

func NewPgxAuditRepository(base *postgres.PgxRepository) audit.AuditRepository {
	return &pgxAuditRepository{
		PgxRepository: base,
	}
}

func (r *pgxAuditRepository) Save(ctx context.Context, log *audit.AuditLog) error {
	query := `
		INSERT INTO audit_logs (log_id, action, actor, resource_id, timestamp)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.Pool.Exec(ctx, query, log.LogID, log.Action, log.Actor, log.ResourceID, log.Timestamp)
	return err
}

func (r *pgxAuditRepository) GetAll(ctx context.Context) ([]audit.AuditLog, error) {
	query := `
		SELECT log_id, action, actor, resource_id, timestamp
		FROM audit_logs
		ORDER BY timestamp DESC
	`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []audit.AuditLog
	for rows.Next() {
		var l audit.AuditLog
		err := rows.Scan(&l.LogID, &l.Action, &l.Actor, &l.ResourceID, &l.Timestamp)
		if err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
