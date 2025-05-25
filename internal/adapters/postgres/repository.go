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

type PostgresRepository struct {
	tx *sql.Tx
}

func NewPostgresRepository(tx *sql.Tx) *PostgresRepository {
	return &PostgresRepository{tx: tx}
}

func (r *PostgresRepository) Save(ctx context.Context, aggregate core.Aggregate) error {
	data, err := json.Marshal(aggregate)
	if err != nil {
		return err
	}

	currentVersion := aggregate.Version()
	newVersion := currentVersion + 1

	if currentVersion == 0 {
		// New aggregate: INSERT
		_, err = r.tx.ExecContext(ctx,
			fmt.Sprintf(`INSERT INTO %s (id, data, version, created_at)
                         VALUES ($1, $2, $3, $4)`, aggregate.TableName()),
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
		res, err := r.tx.ExecContext(ctx,
			fmt.Sprintf(`UPDATE %s
                         SET data = $1, version = $2
                         WHERE id = $3 AND version = $4`, aggregate.TableName()),
			data,
			newVersion, // e.g., version 2
			uuid.UUID(aggregate.ID()),
			currentVersion, // e.g., expecting version 1
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
	baseVersion := currentVersion

	for i, event := range events {
		eventData, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		version := baseVersion + i + 1

		_, err = r.tx.ExecContext(ctx, `
            INSERT INTO events (aggregate_id, aggregate_name, event_type, payload, occurred_at, version)
            VALUES ($1, $2, $3, $4, $5, $6)
        `,
			uuid.UUID(aggregate.ID()),
			aggregate.Name(),
			event.Type(),
			eventData,
			event.OccurredAt(),
			version,
		)
		if err != nil {
			return fmt.Errorf("failed to insert event: %w", err)
		}

		_, err = r.tx.ExecContext(ctx, `
            INSERT INTO outbox (aggregate_id, aggregate_name, event_type, payload, occurred_at, status)
            VALUES ($1, $2, $3, $4, $5, $6)
        `,
			uuid.UUID(aggregate.ID()),
			aggregate.Name(),
			event.Type(),
			eventData,
			time.Now(),
			"pending",
		)
		if err != nil {
			return fmt.Errorf("failed to insert outbox entry: %w", err)
		}
	}

	return nil
}

func (r *PostgresRepository) Load(ctx context.Context, id core.AggregateId, agg core.Aggregate) error {
	var data []byte
	var version int
	var createdAt time.Time

	err := r.tx.QueryRowContext(ctx,
		fmt.Sprintf(`SELECT data, version, created_at FROM %s WHERE id = $1`, agg.TableName()),
		id,
	).Scan(&data, &version, &createdAt)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, agg)
	if err != nil {
		return err
	}

	agg.SetVersion(version)
	agg.SetCreatedAt(createdAt)

	return nil
}
