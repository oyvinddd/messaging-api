CREATE SCHEMA IF NOT EXISTS push;

CREATE TYPE push.platform AS ENUM (
    'apns',
    'fcm'
);

CREATE TABLE push.device_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    account_id UUID NOT NULL,

    token TEXT NOT NULL,
    platform push.platform NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (token)
);

CREATE INDEX idx_device_tokens_account_id
    ON push.device_tokens(account_id);

