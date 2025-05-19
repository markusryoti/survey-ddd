package core

import "context"

type Repository[T Aggregate] interface {
	Save(ctx context.Context, aggregate T) error
	Get(ctx context.Context, id AggregateId) (T, error)
}
