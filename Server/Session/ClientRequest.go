package Session

type RequestType int

const (
	LoginRequest RequestType = iota
	RegisterRequest
	LogoutRequest
	StartGameRequest
	EndGameRequest
	StopWaitingRequest
	ExitRequest
)

var RequestTypeMapper = map[string]RequestType{
	"login":    LoginRequest,
	"register": RegisterRequest,
	"logout":   LogoutRequest,
	"start":    StartGameRequest,
	"end":      EndGameRequest,
	"exit":     ExitRequest,
}

type Request struct {
	Type    RequestType
	Message string
}

type RequestHandler func(*Session, Request) //Todo Decision

func NewRequest(requestType RequestType, message string) *Request {
	return &Request{requestType, message}
}
