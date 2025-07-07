package Communication

type RequestType int

type Request struct {
	Type    RequestType
	Message string
}

type ReplyType int

type Reply struct {
	Type    ReplyType
	Message string
}

const (
	ExitReply ReplyType = iota
	SystemReply
	GameReply
)
