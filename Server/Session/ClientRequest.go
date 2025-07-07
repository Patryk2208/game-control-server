package Session

import "Server/Communication"

const (
	LoginRequest Communication.RequestType = iota
	RegisterRequest
	LogoutRequest
	StartGameRequest
	EndGameRequest
	StopWaitingRequest
	ExitRequest
)

var RequestTypeMapper = map[string]Communication.RequestType{
	"login":    LoginRequest,
	"register": RegisterRequest,
	"logout":   LogoutRequest,
	"start":    StartGameRequest,
	"end":      EndGameRequest,
	"exit":     ExitRequest,
}

type RequestHandler func(*Session, Communication.Request) //Todo Decision

func NewRequest(requestType Communication.RequestType, message string) *Communication.Request {
	return &Communication.Request{Type: requestType, Message: message}
}
