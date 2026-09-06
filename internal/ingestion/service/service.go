package service

import (
	"context"
	"errors"
	"time"

	"github.com/TimurGaliev44/notifyx/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

type notificationPublisher interface {
	Publish(ctx context.Context, notification *domain.Notification) error
}

type notificationRepository interface {
	Insert(ctx context.Context, notification *domain.Notification) error
	Update(ctx context.Context, notification *domain.Notification) error
	GetByIdempotencyKey(ctx context.Context, key string) (*domain.Notification, error)
}

type Service struct {
	repo notificationRepository
	pub  notificationPublisher
}

func New(repo notificationRepository, pub notificationPublisher) *Service {
	return &Service{repo: repo, pub: pub}
}

func (s *Service) Create(ctx context.Context, n *domain.Notification) (bool, error) {
	var lastErr error
	if err := s.repo.Insert(ctx, n); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			existing, err := s.repo.GetByIdempotencyKey(ctx, n.IdempotencyKey)
			if err != nil {
				return false, err
			}
			*n = *existing
			return false, nil
		}
		return false, err
	}

	for i := range 5 {
		if err := s.pub.Publish(ctx, n); err != nil {
			lastErr = err
			select {
			case <-ctx.Done():
				return false, ctx.Err()
			case <-time.After(200 * time.Duration(i+1) * time.Millisecond):
			}
			continue
		}
		lastErr = nil
		break
	}
	if lastErr != nil {
		return false, lastErr
	}
	n.Status = domain.StatusQueued
	if err := s.repo.Update(ctx, n); err != nil {
		return false, err
	}
	return true, nil
}
