package postgres

import (
	"context"
	"database/sql"

	"github.com/markusryoti/survey-ddd/internal/core"
)

type UnitOfWork interface {
	Begin(ctx context.Context) (TransactionalUnit, error)
}

type TransactionalUnit interface {
	Commit() error
	Rollback() error
}

type PostgresUnitOfWork[T core.Aggregate] struct {
	db *sql.DB
}

func NewPostgresUnitOfWork[T core.Aggregate](db *sql.DB) *PostgresUnitOfWork[T] {
	return &PostgresUnitOfWork[T]{db: db}
}

func (u *PostgresUnitOfWork[T]) Begin(ctx context.Context) (TransactionalUnit, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &postgresTransactionalUnit{tx: tx}, nil
}

type postgresTransactionalUnit struct {
	tx *sql.Tx
}

// func (u *postgresTransactionalUnit) LoadAggregate(ctx context.Context, id core.AggregateId) (core.AggregateState, error) {
// 	// implement using shared generic repository that uses u.tx
// 	// e.g., load JSON, parse into aggregate, get version etc.
// }

// func (u *postgresTransactionalUnit) SaveAggregate(ctx context.Context, agg core.Aggregate) error {
// 	// persist aggregate and events using u.tx
// 	return nil
// }

// func (u *postgresTransactionalUnit) WithRepo[T core.Aggregate](tableName string) *Repository[T] {
//     return NewRepositoryWithTx[T](u.tx, tableName, newAgg)
// }

func (u *postgresTransactionalUnit) Commit() error {
	return u.tx.Commit()
}

func (u *postgresTransactionalUnit) Rollback() error {
	return u.tx.Rollback()
}
