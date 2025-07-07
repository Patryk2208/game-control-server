package Session

import (
	"Server/Communication"
	"fmt"
)

type PlayingUserConnectionContext struct {
	PlayingContextRequestMapper map[Communication.RequestType]RequestHandler
}

func (context PlayingUserConnectionContext) GetHandler(request *Communication.Request) (RequestHandler, error) {
	handler, exists := context.PlayingContextRequestMapper[request.Type]
	if !exists {
		return nil, fmt.Errorf("no handler found for request type %s", request.Type)
	}
	return handler, nil
}

func NewPlayingContext() UserConnectionContext {
	return PlayingUserConnectionContext{
		PlayingContextRequestMapper: map[Communication.RequestType]RequestHandler{
			EndGameRequest: EndGameRequestHandler,
		},
	}
}
