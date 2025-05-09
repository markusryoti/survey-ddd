package rabbitmq

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/markusryoti/survey-ddd/internal/core"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

type Outbox struct {
	db        *sqlx.DB
	publisher *Publisher
}

func NewOutbox(db *sqlx.DB, publisher *Publisher) *Outbox {
	return &Outbox{
		db:        db,
		publisher: publisher,
	}
}

func (o *Outbox) Publish(events []core.Event) error {
	ctx := context.Background()
	tx, err := o.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for outbox: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO outbox (aggregate_id, type, payload, occurred_at)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare outbox insert statement: %w", err)
	}
	defer stmt.Close()

	for _, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event payload for outbox: %w", err)
		}
		_, err = stmt.ExecContext(ctx, event.AggregateId(), event.Type(), payload, event.Timestamp()) // Aggregate ID might be relevant here
		if err != nil {
			return fmt.Errorf("failed to insert event into outbox: %w", err)
		}
	}

	return tx.Commit()
}

func (o *Outbox) ProcessOutbox() {
	ticker := time.NewTicker(5 * time.Second) // Check every 5 seconds
	defer ticker.Stop()

	for range ticker.C {
		if err := o.processPendingMessages(); err != nil {
			log.Printf("Error processing outbox messages: %v", err)
		}
	}
}

func (o *Outbox) processPendingMessages() error {
	ctx := context.Background()
	tx, err := o.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for processing outbox: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.QueryxContext(ctx, `
		SELECT id, type, payload
		FROM outbox
		ORDER BY occurred_at ASC
		LIMIT 10 -- Process in batches
	`)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // No messages to process
		}
		return fmt.Errorf("failed to query outbox: %w", err)
	}
	defer rows.Close()

	var processedIDs []int64
	var eventsToPublish []core.Event

	for rows.Next() {
		var outboxMessage struct {
			ID      int64           `db:"id"`
			Type    string          `db:"type"`
			Payload json.RawMessage `db:"payload"`
		}
		if err := rows.StructScan(&outboxMessage); err != nil {
			log.Printf("Error scanning outbox message: %v", err)
			continue
		}

		var event core.Event
		switch outboxMessage.Type {
		case "survey-created":
			event = &surveys.SurveyCreated{}
		case "response-submitted":
			event = &surveys.ResponseSubmitted{}
		// Add other event types here
		default:
			log.Printf("Unknown event type in outbox: %s", outboxMessage.Type)
			processedIDs = append(processedIDs, outboxMessage.ID) // Mark as processed to avoid infinite loops
			continue
		}

		if err := json.Unmarshal(outboxMessage.Payload, event); err != nil {
			log.Printf("Error unmarshalling outbox payload (%s): %v", outboxMessage.Type, err)
			processedIDs = append(processedIDs, outboxMessage.ID) // Mark as processed
			continue
		}

		eventsToPublish = append(eventsToPublish, event)
		processedIDs = append(processedIDs, outboxMessage.ID)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating through outbox rows: %w", err)
	}

	if len(eventsToPublish) > 0 {
		if err := o.publisher.Publish(eventsToPublish); err != nil {
			return fmt.Errorf("failed to publish events from outbox: %w", err)
		}

		// Mark processed messages as deleted
		if len(processedIDs) > 0 {
			_, err = tx.ExecContext(ctx, `
				DELETE FROM outbox
				WHERE id = ANY($1)
			`, processedIDs)
			if err != nil {
				return fmt.Errorf("failed to delete processed outbox messages: %w", err)
			}
		}
	}

	return tx.Commit()
}
