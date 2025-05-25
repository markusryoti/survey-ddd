package core

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, aggregate Aggregate) error
	Load(ctx context.Context, id AggregateId, aggregate Aggregate) error
}

type TransactionSignature func(repo Repository) error

type TransactionProvider interface {
	RunTransactional(ctx context.Context, fn TransactionSignature) error
}
