package main

import (
	"fmt"
)

type PlayingUserConnectionContext struct {
	PlayingContextRequestMapper map[RequestType]RequestHandler
}

func (context PlayingUserConnectionContext) GetHandler(request *Request) (RequestHandler, error) {
	handler, exists := context.PlayingContextRequestMapper[request.Type]
	if !exists {
		return nil, fmt.Errorf("no handler found for request type %s", request.Type)
	}
	return handler, nil
}

func NewPlayingContext() UserConnectionContext {
	return PlayingUserConnectionContext{
		PlayingContextRequestMapper: map[RequestType]RequestHandler{
			EndGameRequest: EndGameRequestHandler,
		},
	}
}
