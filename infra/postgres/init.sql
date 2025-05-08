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
    id UUID PRIMARY KEY,
    aggregate_id UUID NOT NULL,
    type TEXT NOT NULL,
    payload JSONB NOT NULL,
    occurred_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS outbox (
    id UUID PRIMARY KEY,
    aggregate_id UUID NOT NULL,
    type TEXT NOT NULL,
    payload JSONB NOT NULL,
    occurred_at TIMESTAMP NOT NULL,
    published BOOLEAN DEFAULT false
);

