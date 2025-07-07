package Session

import (
	"Server/Communication"
	"fmt"
)

type ReplyHandler func(*Session, Communication.Reply)

var ReplyHandlerMapper = map[Communication.ReplyType]ReplyHandler{
	Communication.ExitReply:   ExitReplyHandler,
	Communication.SystemReply: SystemReplyHandler,
	Communication.GameReply:   GameReplyHandler,
}

func ExitReplyHandler(session *Session, reply Communication.Reply) {
	fmt.Println("Exit Reply Handler")
	panic(fmt.Errorf("exit reply"))
}

func SystemReplyHandler(session *Session, reply Communication.Reply) {
	fmt.Println("System Reply Handler")
	err := session.ClientConn.WriteMessage(1, []byte(reply.Message))
	if err != nil {
		panic(err)
	}
}

func GameReplyHandler(session *Session, reply Communication.Reply) {

}
