package core

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventId uuid.UUID

func (e EventId) String() string {
	return uuid.UUID(e).String()
}

func NewEventId() EventId {
	return EventId(uuid.New())
}

type DomainEvent interface {
	AggregateId() AggregateId
	Type() string
	Timestamp() time.Time
}

func NewEvent(domainEvent DomainEvent) (Event, error) {
	var event Event
	b, err := json.Marshal(domainEvent)
	if err != nil {
		return event, err
	}

	return Event{
		ID:          NewEventId(),
		AggregateID: domainEvent.AggregateId(),
		Type:        domainEvent.Type(),
		Payload:     b,
		OccurredAt:  domainEvent.Timestamp(),
	}, nil
}

type Event struct {
	ID          EventId
	AggregateID AggregateId
	Type        string
	Payload     []byte
	OccurredAt  time.Time
}
