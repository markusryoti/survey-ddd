package core

import "context"

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
