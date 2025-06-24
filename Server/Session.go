package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Session struct {
	Context UserConnectionContext
	//TODO Player Info
	ServerConn *websocket.Conn
	ReplyQueue chan Reply
}

func NewSession() *Session {
	return &Session{
		Context: NewNormalContext(),
	}
}

func (s *Session) Reader() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from ", r)
		}
		s.ReplyQueue <- Reply{ExitReply, "Exit"}
		fmt.Println("Reader Exit")
		return
	}()
	for {
		msgType, msg, err := s.ServerConn.ReadMessage()
		if err != nil {
			panic(err)
		}
		req := CreateRequest(msgType, string(msg))
		handler, err := s.Context.GetHandler(req)
		if err != nil {
			panic(err)
		}
		handler(s, *req)
		fmt.Printf("Message type: %d, %s\n", msgType, string(msg))
	}
}

func (s *Session) Writer() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from ", r)
		}
		fmt.Println("Writer Exit")
		return
	}()
	for reply := range s.ReplyQueue {
		handler, exists := ReplyHandlerMapper[reply.Type]
		if !exists {
			continue
		}
		handler(s, reply)
	}
}
