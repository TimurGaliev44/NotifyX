package domain

import (
	"time"
	"uuid"
)

type Channel string

const (
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusQueued    Status = "queued"
	StatusSent      Status = "sent"
	StatusDelivered Status = "delivered"
	StatusFailed    Status = "failed"
)

type Notification struct {
	ID             uuid.UUID
	IdempotencyKey string
	Channel        Channel
	Recipient      string
	TemplateID     string
	Payload        map[string]any
	Status         Status
	Priority       int
	ScheduledAt    *time.Time
	Source         string
	TraceID        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
