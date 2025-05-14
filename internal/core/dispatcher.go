package core

import (
	"context"
	"sync"
)

type EventHandler func(context.Context, DomainEvent) error

type EventDispatcher struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

func (d *EventDispatcher) Register(eventType string, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

func (d *EventDispatcher) Dispatch(ctx context.Context, events ...DomainEvent) error {
	for _, evt := range events {
		d.mu.RLock()
		handlers := d.handlers[evt.Type()]
		d.mu.RUnlock()

		for _, handler := range handlers {
			if err := handler(ctx, evt); err != nil {
				return err
			}
		}
	}
	return nil
}
