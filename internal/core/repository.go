package core

import "context"

type Repository[T Aggregate] interface {
	Save(ctx context.Context, aggregate T) error
	SaveWithTx(ctx context.Context, tx Transaction, aggregate T) error
	Load(ctx context.Context, id AggregateId, aggregate T) error
	LoadWithTx(ctx context.Context, tx Transaction, id AggregateId, aggregate T) error
}

type Transaction interface {
}

type TransactionalSignature func(tx Transaction) error

type TransactionProvider interface {
	RunTransactional(ctx context.Context, fn TransactionalSignature) error
}
