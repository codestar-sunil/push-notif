-- 0001_create_devices.up.sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE devices (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  app_id          UUID NOT NULL,
  installation_id TEXT NOT NULL,
  push_token      TEXT NOT NULL,
  platform        TEXT NOT NULL CHECK (platform IN ('ios','android','web')),
  user_id         TEXT,
  last_seen       TIMESTAMPTZ NOT NULL DEFAULT now(),
  status          TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','stale','unregistered')),
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (app_id, installation_id)
);

CREATE INDEX idx_devices_app_status ON devices (app_id, status);

