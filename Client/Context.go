package main

import "errors"

type ClientContext int

const (
	Normal ClientContext = iota
	Authenticated
	Playing
)

type Context interface {
	GetHandler(command Command) (CommandHandler, error)
	GetPrompt() string
}

type NormalContext struct {
	CommandHandlers map[string]CommandHandler
}

func (context NormalContext) GetHandler(command Command) (CommandHandler, error) {
	contextHandler, exists := context.CommandHandlers[command.Name]
	if !exists {
		return nil, errors.New("command not found")
	}
	return contextHandler, nil
}

func (context NormalContext) GetPrompt() string { return "user" }

type AuthenticatedContext struct {
	CommandHandlers map[string]CommandHandler
}

func (context AuthenticatedContext) GetHandler(command Command) (CommandHandler, error) {
	contextHandler, exists := context.CommandHandlers[command.Name]
	if !exists {
		return nil, errors.New("command not found")
	}
	return contextHandler, nil
}

func (context AuthenticatedContext) GetPrompt() string { return "auth" }

type PlayingContext struct {
	CommandHandlers map[string]CommandHandler
}

func (context PlayingContext) GetHandler(command Command) (CommandHandler, error) {
	contextHandler, exists := context.CommandHandlers[command.Name]
	if !exists {
		return nil, errors.New("command not found")
	}
	return contextHandler, nil
}

func (context PlayingContext) GetPrompt() string { return "playing" }

func NewNormalContext() Context {
	return NormalContext{
		CommandHandlers: map[string]CommandHandler{
			"login":    LoginCommandHandler,
			"register": RegisterCommandHandler,
			"help":     NormalHelpCommandHandler,
			"exit":     ExitCommandHandler,
		},
	}
}

func NewAuthenticatedContext() Context {
	return AuthenticatedContext{
		CommandHandlers: map[string]CommandHandler{
			"logout": LogoutCommandHandler,
			"play":   PlayCommandHandler,
			"help":   AuthenticatedHelpCommandHandler,
			"exit":   ExitCommandHandler,
		},
	}
}

func NewPlayingContext() Context {
	return PlayingContext{
		CommandHandlers: map[string]CommandHandler{
			"quit": ExitCommandHandler,
		},
	}
}
