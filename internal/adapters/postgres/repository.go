package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/markusryoti/survey-ddd/internal/core"
)

type PostgresRepository[T core.Aggregate] struct {
	db        *sql.DB
	tableName string
}

func NewPostgresRepository[T core.Aggregate](db *sql.DB, tableName string) *PostgresRepository[T] {
	return &PostgresRepository[T]{db: db, tableName: tableName}
}

func (r *PostgresRepository[T]) Save(ctx context.Context, aggregate T) error {
	var err error

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func(cause error) {
		if cause != nil {
			tx.Rollback()
		}
	}(err)

	data, err := json.Marshal(aggregate)
	if err != nil {
		return err
	}

	if aggregate.Version() == 0 {
		// New aggregate: INSERT
		_, err = tx.ExecContext(ctx,
			fmt.Sprintf(`INSERT INTO %s (id, data, version, created_at)
                         VALUES ($1, $2, $3, $4)`, r.tableName),
			uuid.UUID(aggregate.ID()),
			data,
			1, // First version number is 1
			aggregate.CreatedAt(),
		)
		if err != nil {
			return fmt.Errorf("insert failed: %w", err)
		}
	} else {
		// Existing aggregate: UPDATE with OCC
		res, err := tx.ExecContext(ctx,
			fmt.Sprintf(`UPDATE %s
                         SET data = $1, version = $2
                         WHERE id = $3 AND version = $4`, r.tableName),
			data,
			aggregate.Version(), // e.g., version 2
			uuid.UUID(aggregate.ID()),
			aggregate.Version()-1, // e.g., expecting version 1 in DB
		)
		if err != nil {
			return fmt.Errorf("update failed: %w", err)
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("rows affected error: %w", err)
		}
		if rowsAffected == 0 {
			return errors.New("optimistic concurrency conflict: aggregate has been modified")
		}
	}

	events := aggregate.GetUncommittedEvents()
	nextVersion := aggregate.Version() - len(events) + 1

	for _, event := range events {
		// Serialize event
		eventData, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		// Insert into event store
		_, err = tx.ExecContext(ctx, `
            INSERT INTO events (aggregate_id, type, payload, occurred_at, version)
            VALUES ($1, $2, $3, $4, $5)
        `,
			uuid.UUID(aggregate.ID()),
			event.Type(),
			eventData,
			time.Now(),
			nextVersion,
		)
		if err != nil {
			return fmt.Errorf("failed to insert event: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
            INSERT INTO outbox (aggregate_id, type, payload, occurred_at, status)
            VALUES ($1, $2, $3, $4, $5)
        `,
			uuid.UUID(aggregate.ID()),
			event.Type(),
			eventData,
			time.Now(),
			"pending", // or whatever your initial status is
		)
		if err != nil {
			return fmt.Errorf("failed to insert outbox entry: %w", err)
		}

		nextVersion++
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository[T]) Load(ctx context.Context, id core.AggregateId, agg core.Aggregate) error {
	var data []byte
	var version int
	var createdAt time.Time

	err := r.db.QueryRowContext(ctx,
		fmt.Sprintf(`SELECT data, version, created_at FROM %s WHERE id = $1`, r.tableName),
		id,
	).Scan(&data, &version, &createdAt)

	err = json.Unmarshal(data, agg)
	if err != nil {
		return err
	}

	agg.SetVersion(version)
	agg.SetCreatedAt(createdAt)

	return nil
}
