package Session

import (
	"Server/Communication"
	"fmt"
)

type WaitingUserConnectionContext struct {
	WaitingContextRequestMapper map[Communication.RequestType]RequestHandler
}

func (context WaitingUserConnectionContext) GetHandler(request *Communication.Request) (RequestHandler, error) {
	handler, exists := context.WaitingContextRequestMapper[request.Type]
	if !exists {
		return nil, fmt.Errorf("no handler found for request type %d", request.Type)
	}
	return handler, nil
}

func NewWaitingContext() UserConnectionContext {
	return WaitingUserConnectionContext{
		WaitingContextRequestMapper: map[Communication.RequestType]RequestHandler{
			StopWaitingRequest: StopWaitingRequestHandler,
			LogoutRequest:      StopWaitingAndLogoutRequestHandler,
			ExitRequest:        StopWaitingAndExitRequestHandler,
		},
	}
}
