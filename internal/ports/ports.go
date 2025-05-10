package ports

import (
	"context"

	"github.com/markusryoti/survey-ddd/internal/core"
)

type Repository[T core.Aggregate] interface {
	Save(ctx context.Context, aggregate T) error
	Get(ctx context.Context, id core.AggregateId) (T, error)
}
