package main

import (
	"fmt"
)

type AuthenticatedUserConnectionContext struct {
	AuthenticatedContextRequestMapper map[RequestType]RequestHandler
}

func (context AuthenticatedUserConnectionContext) GetHandler(request *Request) (RequestHandler, error) {
	handler, exists := context.AuthenticatedContextRequestMapper[request.Type]
	if !exists {
		return nil, fmt.Errorf("no handler found for request type %s", request.Type)
	}
	return handler, nil
}

func NewAuthenticatedContext() UserConnectionContext {
	return AuthenticatedUserConnectionContext{
		AuthenticatedContextRequestMapper: map[RequestType]RequestHandler{
			LogoutRequest:    LogoutRequestHandler,
			StartGameRequest: StartGameRequestHandler,
			ExitRequest:      ExitWithLogoutRequestHandler,
		},
	}
}
