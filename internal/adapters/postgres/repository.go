package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

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
			aggregate.ID(),
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
			aggregate.ID(),
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
