package postgres

import (
	"context"
	"database/sql"
)

type PostgresTx struct {
	db *sql.DB
	tx *sql.Tx
}

func (p *PostgresTx) Begin(ctx context.Context) error {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	p.tx = tx

	return nil
}

func (p *PostgresTx) Commit() error {
	return p.tx.Commit()
}

func (p *PostgresTx) Rollback() error {
	return p.tx.Rollback()
}

func (p *PostgresTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return p.tx.ExecContext(ctx, query, args)
}
