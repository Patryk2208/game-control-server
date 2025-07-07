package Session

import (
	"Server/Communication"
	"fmt"
)

type NormalUserConnectionContext struct {
	NormalContextRequestMapper map[Communication.RequestType]RequestHandler
}

func (context NormalUserConnectionContext) GetHandler(request *Communication.Request) (RequestHandler, error) {
	handler, exists := context.NormalContextRequestMapper[request.Type]
	if !exists {
		return nil, fmt.Errorf("no handler found for request type %s", string(request.Type))
	}
	return handler, nil
}

func NewNormalContext() UserConnectionContext {
	return NormalUserConnectionContext{
		NormalContextRequestMapper: map[Communication.RequestType]RequestHandler{
			LoginRequest:    LoginRequestHandler,
			RegisterRequest: RegisterRequestHandler,
			ExitRequest:     ExitRequestHandler,
		},
	}
}
