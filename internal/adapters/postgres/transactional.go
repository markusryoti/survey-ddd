package postgres

import (
	"context"
	"database/sql"

	"github.com/markusryoti/survey-ddd/internal/core"
)

func NewPostgresTransactionalProvider(db *sql.DB) *PostgresTransactionalProvider {
	return &PostgresTransactionalProvider{
		db: db,
	}
}

type PostgresTransactionalProvider struct {
	db *sql.DB
}

type PostgresTx struct {
	tx *sql.Tx
}

func (p *PostgresTransactionalProvider) RunTransactional(ctx context.Context, fn core.TransactionalSignature) error {
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

	t := &PostgresTx{tx: tx}

	err = fn(t)
	if err != nil {
		return err
	}

	return tx.Commit()
}
