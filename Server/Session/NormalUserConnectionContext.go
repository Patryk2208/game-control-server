package Session

import (
	"fmt"
)

type NormalUserConnectionContext struct {
	NormalContextRequestMapper map[RequestType]RequestHandler
}

func (context NormalUserConnectionContext) GetHandler(request *Request) (RequestHandler, error) {
	handler, exists := context.NormalContextRequestMapper[request.Type]
	if !exists {
		return nil, fmt.Errorf("no handler found for request type %s", string(request.Type))
	}
	return handler, nil
}

func NewNormalContext() UserConnectionContext {
	return NormalUserConnectionContext{
		NormalContextRequestMapper: map[RequestType]RequestHandler{
			LoginRequest:    LoginRequestHandler,
			RegisterRequest: RegisterRequestHandler,
			ExitRequest:     ExitRequestHandler,
		},
	}
}
