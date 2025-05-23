package core

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AggregateId uuid.UUID

func (a AggregateId) String() string {
	return uuid.UUID(a).String()
}

type Aggregate interface {
	ID() AggregateId
	GetUncommittedEvents() []DomainEvent
	AddDomainEvent(DomainEvent)
	ClearUncommittedEvents()
	SetVersion(int)
	SetCreatedAt(time.Time)
	Version() int
	CreatedAt() time.Time
	Name() string
	TableName() string
}

type BaseAggregate struct {
	version           int
	createdAt         time.Time
	uncommittedEvents []DomainEvent
}

func (b BaseAggregate) GetUncommittedEvents() []DomainEvent {
	return b.uncommittedEvents
}

func (b *BaseAggregate) AddDomainEvent(e DomainEvent) {
	b.uncommittedEvents = append(b.uncommittedEvents, e)
}

func (b *BaseAggregate) ClearUncommittedEvents() {
	b.uncommittedEvents = make([]DomainEvent, 0)
}

func (b *BaseAggregate) SetVersion(version int) {
	b.version = version
}

func (b *BaseAggregate) SetCreatedAt(t time.Time) {
	b.createdAt = t
}

func (b BaseAggregate) Version() int {
	return b.version
}

func (b BaseAggregate) CreatedAt() time.Time {
	return b.createdAt
}

func NewAggregateId() AggregateId {
	return AggregateId(uuid.New())
}

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
