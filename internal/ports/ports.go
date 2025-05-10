package ports

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/markusryoti/survey-ddd/internal/core"
)

type Repository[T core.Aggregate] interface {
	Save(ctx context.Context, aggregate T) error
	Get(ctx context.Context, id core.AggregateId) (T, error)
}

type EventStore interface {
	AppendToStream(ctx context.Context, tx *sqlx.Tx, aggregateID core.AggregateId, events []core.DomainEvent, expectedVersion int) error
	LoadStream(ctx context.Context, aggregateID core.AggregateId) ([]core.DomainEvent, error)
}

type EventPublisher interface {
	Publish(events []core.DomainEvent) error
	Close() error
}

type Outbox interface {
	Publish(events []core.DomainEvent) error
	ProcessOutbox()
}
