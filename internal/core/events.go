package core

import (
	"time"
)

type DomainEvent interface {
	AggregateId() AggregateId
	Type() string
	OccurredAt() time.Time
}
