package core

import (
	"time"

	"github.com/google/uuid"
)

type AggregateId uuid.UUID

type Aggregate interface {
	ID() AggregateId
	GetUncommittedEvents() []Event
	ClearUncommittedEvents()
	SetVersion(int)
	Version() int
	CreatedAt() time.Time
}

type Event interface {
	AggregateId() AggregateId
	Type() string
	Timestamp() time.Time
}
