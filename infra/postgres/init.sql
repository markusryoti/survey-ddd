CREATE TABLE IF NOT EXISTS question_options (
    id UUID PRIMARY KEY,
    value TEXT NOT NULL,
    question_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_question
        FOREIGN KEY(question_id)
            REFERENCES questions(id)
);

CREATE TABLE IF NOT EXISTS questions (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    survey_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_survey
        FOREIGN KEY(survey_id)
            REFERENCES surveys(id)
);

CREATE TABLE IF NOT EXISTS surveys (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    aggregate_id UUID NOT NULL,
    type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version INTEGER NOT NULL,
    CONSTRAINT unique_aggregate_version UNIQUE (aggregate_id, version)
);

CREATE INDEX idx_events_aggregate_id ON events (aggregate_id);
CREATE INDEX idx_events_type ON events (type);
CREATE INDEX idx_events_occurred_at ON events (occurred_at);

CREATE TABLE IF NOT EXISTS outbox (
    id BIGSERIAL PRIMARY KEY,
    aggregate_id UUID,
    type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_outbox_occurred_at ON outbox (occurred_at);