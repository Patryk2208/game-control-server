package main

import (
	"fmt"
	"strings"
)

type ReplyType int

const (
	SystemReply ReplyType = iota
	GameReply
)

var ReplyTypeMapper = map[int]ReplyType{
	1: SystemReply,
	2: GameReply,
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
	GameReply:   GameReplyHandler,
}

func SystemReplyHandler(session *Session, reply Reply) {
	session.ReplyChannel <- reply
	split := strings.Split(reply.Message, " ")
	if len(split) > 1 && split[1] == "exit" {
		panic(fmt.Errorf("received exit reply"))
	}
}

func GameReplyHandler(session *Session, reply Reply) {
	//todo game
}
