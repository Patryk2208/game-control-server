package main

import (
	"Client/GameplayManager"
	"fmt"
	"strings"
)

type ReplyType int

const (
	SystemReply ReplyType = iota
)

var ReplyTypeMapper = map[int]ReplyType{
	1: SystemReply,
}

type Reply struct {
	Type    ReplyType
	Message string
}

func NewReply(t ReplyType, message string) Reply {
	return Reply{Type: t, Message: message}
}

type ReplyHandler func(session *Session, reply Reply)

var ReplyHandlerMapper = map[ReplyType]ReplyHandler{
	SystemReply: SystemReplyHandler,
}

func SystemReplyHandler(session *Session, reply Reply) {
	if len(reply.Message) <= 0 {
		return
	} else if reply.Message[0] == 'G' {
		GameplayManager.RunGameplay(reply.Message[2:])
		session.OperationComplete.L.Lock()
		session.OperationComplete.Signal()
		session.OperationComplete.L.Unlock()
	} else if reply.Message[0] == 'T' || reply.Message[0] == 'F' {
		session.ReplyChannel <- reply
		split := strings.Split(reply.Message, " ")
		if len(split) > 1 && split[1] == "exit" {
			panic(fmt.Errorf("received exit reply"))
		}
	}
}
