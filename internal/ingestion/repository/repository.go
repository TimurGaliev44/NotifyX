package repository

import (
	"context"
	"fmt"

	"github.com/TimurGaliev44/notifyx/internal/domain"
	"github.com/TimurGaliev44/notifyx/internal/pkg/postgres"
)

type Repository struct {
	db *postgres.Database
}

func New(db *postgres.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Insert(ctx context.Context, n *domain.Notification) error {
	const query = `
        INSERT INTO notifications
        (idempotency_key, channel, recipient, template_id, payload, priority, scheduled_at, source, trace_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at, updated_at, status`
	return r.db.Pool.QueryRow(ctx, query,
		n.IdempotencyKey, n.Channel, n.Recipient, n.TemplateID,
		n.Payload, n.Priority, n.ScheduledAt, n.Source, n.TraceID,
	).Scan(&n.ID, &n.CreatedAt, &n.UpdatedAt, &n.Status)
}

func (r *Repository) Update(ctx context.Context, notification domain.Notification) error {
	sqlQuery := `UPDATE notifications SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Pool.Exec(ctx, sqlQuery, notification.Status, notification.UpdatedAt, notification.ID)
	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}
	return nil
}
