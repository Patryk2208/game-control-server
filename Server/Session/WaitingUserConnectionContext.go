package Session

import (
	"fmt"
)

type WaitingUserConnectionContext struct {
	WaitingContextRequestMapper map[RequestType]RequestHandler
}

func (context WaitingUserConnectionContext) GetHandler(request *Request) (RequestHandler, error) {
	handler, exists := context.WaitingContextRequestMapper[request.Type]
	if !exists {
		return nil, fmt.Errorf("no handler found for request type %d", request.Type)
	}
	return handler, nil
}

func NewWaitingContext() UserConnectionContext {
	return PlayingUserConnectionContext{
		PlayingContextRequestMapper: map[RequestType]RequestHandler{
			StopWaitingRequest: StopWaitingRequestHandler,
			LogoutRequest:      StopWaitingAndLogoutRequestHandler,
			ExitRequest:        StopWaitingAndExitRequestHandler,
		},
	}
}
