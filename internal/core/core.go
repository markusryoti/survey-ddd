package core

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AggregateId uuid.UUID
type EventId uuid.UUID

type Aggregate interface {
	ID() AggregateId
	GetUncommittedEvents() []DomainEvent
	ClearUncommittedEvents()
	SetVersion(int)
	Version() int
	CreatedAt() time.Time
}

func NewAggregateId() AggregateId {
	return AggregateId(uuid.New())
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

type Command interface {
	Type() string
}

type CommandBus interface {
	Register(commandType string, handler CommandHandler)
	Dispatch(command Command) error
}

type CommandHandler interface {
	Handle(command Command) error
}
