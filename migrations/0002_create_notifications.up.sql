CREATE TABLE notifications (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id        UUID NOT NULL,
    title         TEXT NOT NULL,
    body          TEXT NOT NULL,
    device_id     UUID,
    segment_id    UUID,
    broadcast     BOOLEAN NOT NULL DEFAULT false,
    schedule_time TIMESTAMPTZ,
    status        TEXT NOT NULL DEFAULT 'pending'
                      CHECK (status IN ('pending','processing','sent','failed')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT exactly_one_target CHECK (...)
);
CREATE INDEX idx_notifications_app_status ON notifications (app_id, status);
CREATE INDEX idx_notifications_schedule ON notifications (schedule_time)
    WHERE status = 'pending' AND schedule_time IS NOT NULL;