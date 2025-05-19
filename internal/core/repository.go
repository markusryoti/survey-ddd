package core

import "context"

type Repository[T Aggregate] interface {
	Save(ctx context.Context, aggregate T) error
	SaveWithTx(ctx context.Context, tx Transactional, aggregate T) error
	Load(ctx context.Context, id AggregateId, aggregate T) error
	LoadWithTx(ctx context.Context, tx Transactional, id AggregateId, aggregate T) error
}

type Transactional interface {
	Begin(ctx context.Context) error
	Commit() error
	Rollback() error
}
