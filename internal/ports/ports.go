package ports

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/markusryoti/survey-ddd/internal/core"
)

type Repository[T core.Aggregate] interface {
	Get(ctx context.Context, id core.AggregateId) (T, error)
	Save(ctx context.Context, aggregate T, events []core.Event) error
}

type EventStore interface {
	AppendToStream(ctx context.Context, tx *sqlx.Tx, aggregateID core.AggregateId, events []core.Event, expectedVersion int) error
	LoadStream(ctx context.Context, aggregateID core.AggregateId) ([]core.Event, error)
}

type EventPublisher interface {
	Publish(events []core.Event) error
	Close() error
}

type Outbox interface {
	Publish(events []core.Event) error
	ProcessOutbox()
}
