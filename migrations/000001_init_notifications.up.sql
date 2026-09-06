CREATE TYPE notification_channel AS ENUM ('email', 'sms');
CREATE TYPE notification_status AS ENUM ('pending', 'queued', 'sent', 'delivered', 'failed');

CREATE TABLE notifications (
                               id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                               idempotency_key  TEXT UNIQUE,
                               channel          notification_channel NOT NULL,
                               recipient        TEXT NOT NULL,
                               template_id      TEXT NOT NULL,
                               payload          JSONB NOT NULL DEFAULT '{}',
                               status           notification_status NOT NULL DEFAULT 'pending',
                               priority         SMALLINT NOT NULL DEFAULT 0,
                               scheduled_at     TIMESTAMPTZ,
                               source           TEXT NOT NULL,
                               trace_id         TEXT,
                               created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
                               updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_notifications_status_created ON notifications (status, created_at);