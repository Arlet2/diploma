CREATE TYPE PushStatus as ENUM (
    'ON_DELIVERY',
    'DELIVERED',
    'NACKED'
);

CREATE TABLE pushes (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    device_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    status PushStatus NOT NULL
);