CREATE SCHEMA IF NOT EXISTS push;

CREATE TYPE push.push_platform AS ENUM (
    'apns',
    'fcm'
);

CREATE TABLE push.tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    account_id UUID NOT NULL,

    device_token TEXT NOT NULL,
    platform push_platform NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (token)
);
