package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/markusryoti/survey-ddd/internal/core"
)

type PostgresRepository[T core.Aggregate] struct {
	db         *sql.DB
	tableName  string
	newAggFunc func() T
}

func NewPostgresRepository[T core.Aggregate](db *sql.DB, tableName string, newAggFunc func() T) *PostgresRepository[T] {
	return &PostgresRepository[T]{db: db, tableName: tableName, newAggFunc: newAggFunc}
}

func (r *PostgresRepository[T]) Save(ctx context.Context, aggregate T) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	data, err := json.Marshal(aggregate)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx,
		fmt.Sprintf(`INSERT INTO %s (id, data, version, created_at) 
                    VALUES ($1, $2, $3, $4) 
                    ON CONFLICT (id) 
                    DO UPDATE SET data = $2, version = $3`, r.tableName),
		aggregate.ID(), data, aggregate.Version(), aggregate.CreatedAt(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// save uncommitted events
	for _, domainEvent := range aggregate.GetUncommittedEvents() {
		event, err := core.NewEvent(domainEvent)
		if err != nil {
			return err
		}

		eventId := uuid.New()
		_, err = tx.ExecContext(ctx,
			`INSERT INTO events (id, aggregate_id, type, payload, occurred_at) VALUES ($1, $2, $3, $4, $5)`,
			eventId, domainEvent.AggregateId(), event.Type, event.Payload, event.OccurredAt)
		if err != nil {
			tx.Rollback()
			return err
		}

		// insert into outbox
		_, err = tx.ExecContext(ctx,
			`INSERT INTO outbox (id, aggregate_id, type, payload, occurred_at, published) VALUES ($1, $2, $3, $4, $5, false)`,
			event.ID, event.AggregateID, event.Type, event.Payload, event.OccurredAt)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	aggregate.ClearUncommittedEvents()
	return nil
}

func (r *PostgresRepository[T]) Load(ctx context.Context, id core.AggregateId) (T, error) {
	agg := r.newAggFunc()
	var data []byte
	err := r.db.QueryRowContext(ctx,
		fmt.Sprintf(`SELECT data FROM %s WHERE id = $1`, r.tableName),
		id,
	).Scan(&data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return agg, fmt.Errorf("aggregate not found")
		}
		return agg, err
	}

	err = json.Unmarshal(data, agg)
	if err != nil {
		return agg, err
	}
	return agg, nil
}
