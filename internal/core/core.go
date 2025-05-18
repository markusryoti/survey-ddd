package core

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AggregateId uuid.UUID

func (id AggregateId) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *AggregateId) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := uuid.Parse(s)
	if err != nil {
		return err
	}
	*id = AggregateId(parsed)
	return nil
}

func (id AggregateId) Value() (driver.Value, error) {
	return uuid.UUID(id).String(), nil
}

func (id *AggregateId) Scan(value interface{}) error {
	if value == nil {
		return errors.New("null UUID")
	}

	switch v := value.(type) {
	case []byte:
		parsed, err := uuid.ParseBytes(v)
		if err != nil {
			return err
		}
		*id = AggregateId(parsed)
		return nil
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		*id = AggregateId(parsed)
		return nil
	default:
		return errors.New("invalid type for UUID")
	}
}

func (a AggregateId) String() string {
	return uuid.UUID(a).String()
}

type EventId uuid.UUID

func (e EventId) String() string {
	return uuid.UUID(e).String()
}

type Aggregate interface {
	ID() AggregateId
	GetUncommittedEvents() []DomainEvent
	ClearUncommittedEvents()
	SetVersion(int)
	SetCreatedAt(time.Time)
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
	Handle(ctx context.Context, command Command) error
}
