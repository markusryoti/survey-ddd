package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/markusryoti/survey-ddd/internal/core"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
	"github.com/markusryoti/survey-ddd/internal/ports"
)

type SurveyRepository struct {
	db         *sqlx.DB
	eventStore ports.EventStore
}

func NewSurveyRepository(db *sqlx.DB, eventStore ports.EventStore) *SurveyRepository {
	return &SurveyRepository{db: db, eventStore: eventStore}
}

func (r *SurveyRepository) Save(ctx context.Context, aggregate *surveys.Survey, events []core.Event) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Save the aggregate state (or its representation)
	_, err = tx.ExecContext(ctx, `
		INSERT INTO surveys (id, title, created_at, version)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET title = $2, version = $4;
	`, aggregate.ID, aggregate.Title, aggregate.CreatedAt(), aggregate.Version())
	if err != nil {
		return fmt.Errorf("failed to save aggregate: %w", err)
	}

	// 2. Save the domain events to the Event Store within the same transaction
	err = r.eventStore.AppendToStream(ctx, tx, aggregate.ID(), events, aggregate.Version()-len(events))
	if err != nil {
		return fmt.Errorf("failed to save events: %w", err)
	}

	// 3. Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Optionally clear uncommitted events on the aggregate after successful persistence
	aggregate.ClearUncommittedEvents()

	return nil
}

func (r *SurveyRepository) Get(id string) (*surveys.Survey, error) {
	// ... (Implementation to retrieve aggregate state and load events)
	// You would typically load the aggregate state and then replay events
	// from the Event Store to reconstruct the current state.
	return nil, fmt.Errorf("get implementation not yet complete")
}
