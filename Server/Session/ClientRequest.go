package Session

import "Server/Communication"

const (
	LoginRequest Communication.RequestType = iota
	RegisterRequest
	LogoutRequest
	StartGameRequest
	StopWaitingRequest
	ExitRequest
)

var RequestTypeMapper = map[string]Communication.RequestType{
	"login":    LoginRequest,
	"register": RegisterRequest,
	"logout":   LogoutRequest,
	"start":    StartGameRequest,
	"stop":     StopWaitingRequest,
	"exit":     ExitRequest,
}

type RequestHandler func(*Session, Communication.Request)

func NewRequest(requestType Communication.RequestType, message string) *Communication.Request {
	return &Communication.Request{Type: requestType, Message: message}
}
