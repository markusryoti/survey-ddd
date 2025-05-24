CREATE TABLE surveys (
    id UUID PRIMARY KEY,
    data JSONB NOT NULL,
    version INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE survey_responses (
    id UUID PRIMARY KEY,
    data JSONB NOT NULL,
    version INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    aggregate_id UUID NOT NULL,
    aggregate_name VARCHAR(255),
    event_type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version INTEGER NOT NULL,
    CONSTRAINT unique_aggregate_version UNIQUE (aggregate_id, version)
);

CREATE INDEX idx_events_aggregate_id ON events (aggregate_id);
CREATE INDEX idx_events_type ON events (event_type);
CREATE INDEX idx_events_occurred_at ON events (occurred_at);

CREATE TABLE IF NOT EXISTS outbox (
    id BIGSERIAL PRIMARY KEY,
    aggregate_id UUID,
    aggregate_name VARCHAR(255),
    event_type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_outbox_occurred_at ON outbox (occurred_at);