package main

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
	Accepted
	Rejected
	GameCommand
)

type ReplyHandler func(*Session, Reply)

var ReplyHandlerMapper = map[ReplyType]ReplyHandler{
	ExitReply: ExitReplyHandler,
	Accepted:  AcceptedHandler,
	Rejected:  RejectedHandler,
}

func ExitReplyHandler(session *Session, reply Reply) {
	fmt.Println("Exit Reply Handler")
	panic(fmt.Errorf("exit reply"))
}

func AcceptedHandler(session *Session, reply Reply) {
	fmt.Println("Accepted Reply Handler")
	err := session.ServerConn.WriteMessage(int(reply.Type), []byte(reply.Message))
	if err != nil {
		panic(err)
	}
}

func RejectedHandler(session *Session, reply Reply) {
	fmt.Println("Rejected Reply Handler")
	err := session.ServerConn.WriteMessage(int(reply.Type), []byte(reply.Message))
	if err != nil {
		panic(err)
	}
}
