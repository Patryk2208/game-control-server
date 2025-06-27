package Session

import (
	"fmt"
)

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

type ReplyHandler func(*Session, Reply)

var ReplyHandlerMapper = map[ReplyType]ReplyHandler{
	ExitReply:   ExitReplyHandler,
	SystemReply: SystemReplyHandler,
	GameReply:   GameReplyHandler,
}

func ExitReplyHandler(session *Session, reply Reply) {
	fmt.Println("Exit Reply Handler")
	panic(fmt.Errorf("exit reply"))
}

func SystemReplyHandler(session *Session, reply Reply) {
	fmt.Println("System Reply Handler")
	err := session.ClientConn.WriteMessage(1, []byte(reply.Message))
	if err != nil {
		panic(err)
	}
}

func GameReplyHandler(session *Session, reply Reply) {

}
