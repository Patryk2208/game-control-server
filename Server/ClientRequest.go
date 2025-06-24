package main

type RequestType int

const (
	LoginRequest RequestType = iota
	RegisterRequest
	LogoutRequest
	StartGameRequest
	EndGameRequest
	ExitRequest
)

type Request struct {
	Type    RequestType
	Message string
}

type RequestHandler func(*Session, Request) //Todo Decision

func NewRequest(requestType RequestType, message string) *Request {
	return &Request{requestType, message}
}

func CreateRequest(requestType int, message string) *Request {
	return &Request{RequestType(requestType), message}
}
