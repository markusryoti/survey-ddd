package postgres

import (
	"context"
	"database/sql"

	"github.com/markusryoti/survey-ddd/internal/core"
)

type PostgresTransactionalProvider struct {
	db *sql.DB
}

func NewPostgresTransactionalProvider(db *sql.DB) *PostgresTransactionalProvider {
	return &PostgresTransactionalProvider{
		db: db,
	}
}

func (p *PostgresTransactionalProvider) RunTransactional(ctx context.Context, fn core.TransactionSignature) error {
	var err error
	var tx *sql.Tx

	tx, err = p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				panic("rollback failed")
			}
		}
	}()

	t := NewPostgresRepository(tx)

	err = fn(t)
	if err != nil {
		return err
	}

	return tx.Commit()
}
