CREATE TABLE deliveries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id UUID NOT NULL,
    device_id       UUID NOT NULL,
    status          TEXT NOT NULL DEFAULT 'queued'
                        CHECK (status IN ('queued','sent','delivered','failed','bounced')),
    attempt_count   INT NOT NULL DEFAULT 0,
    last_error      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (notification_id, device_id)
);
CREATE INDEX idx_deliveries_notification ON deliveries (notification_id);