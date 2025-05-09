package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/markusryoti/survey-ddd/internal/core"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

type EventStore struct {
	db *sqlx.DB
}

func NewEventStore(db *sqlx.DB) *EventStore {
	return &EventStore{db: db}
}

func (es *EventStore) AppendToStream(ctx context.Context, tx *sqlx.Tx, aggregateID string, events []core.Event, expectedVersion int) error {
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO events (aggregate_id, type, payload, occurred_at, version)
		SELECT $1, $2, $3, $4, (SELECT COALESCE(MAX(version), 0) + 1 FROM events WHERE aggregate_id = $1)
		WHERE NOT EXISTS (SELECT 1 FROM events WHERE aggregate_id = $1 AND version = $5);
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare event insert statement: %w", err)
	}
	defer stmt.Close()

	for _, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event payload: %w", err)
		}
		_, err = stmt.ExecContext(ctx, aggregateID, event.Type(), payload, event.Timestamp(), expectedVersion+1) // Increment expected version for each event
		if err != nil {
			return fmt.Errorf("failed to insert event: %w", err)
		}
		expectedVersion++
	}

	return nil
}

func (es *EventStore) LoadStream(aggregateID string) ([]core.Event, error) {
	var records []EventRecord
	err := es.db.Select(&records, "SELECT id, aggregate_id, type, payload, occurred_at, version FROM events WHERE aggregate_id = $1 ORDER BY version ASC", aggregateID)
	if err != nil {
		return nil, fmt.Errorf("failed to load events for aggregate %s: %w", aggregateID, err)
	}

	var events []core.Event
	for _, record := range records {
		var event core.Event
		switch record.Type {
		case "survey-created":
			event = &surveys.SurveyCreated{}
		case "response-submitted":
			event = &surveys.ResponseSubmitted{}
		// Add other event types here
		default:
			continue
		}
		if err := json.Unmarshal(record.Payload, event); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event payload: %w", err)
		}
		events = append(events, event)
	}
	return events, nil
}

type EventRecord struct {
	ID          int64     `db:"id"`
	AggregateID string    `db:"aggregate_id"`
	Type        string    `db:"type"`
	Payload     []byte    `db:"payload"`
	OccurredAt  time.Time `db:"occurred_at"`
}
