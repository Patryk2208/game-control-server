package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type Session struct {
	Context Context
	//todo player info
	ServerConn        *websocket.Conn
	CommandQueue      chan Command
	ReplyChannel      chan Reply
	OperationComplete *sync.Cond
}

func NewSession() *Session {
	return &Session{
		Context:           NewNormalContext(),
		CommandQueue:      make(chan Command, 100),
		ReplyChannel:      make(chan Reply, 100),
		OperationComplete: sync.NewCond(&sync.Mutex{}),
	}
}

func (s *Session) Close() {
	if s.ServerConn != nil {
		s.ServerConn.Close()
	}
	close(s.CommandQueue)
	close(s.ReplyChannel)
}

func (s *Session) Reader(wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from ", r)
		}
		fmt.Println("Reader Exit")
		wg.Done()
		return
	}()
	for {
		msgType, msg, err := s.ServerConn.ReadMessage()
		if err != nil {
			panic(err)
		}
		replyType, exists := ReplyTypeMapper[msgType]
		if !exists {
			fmt.Println("Unknown reply type:", msgType)
			s.OperationComplete.L.Lock()
			s.OperationComplete.Signal()
			s.OperationComplete.L.Unlock()
			continue
		}
		reply := NewReply(replyType, string(msg))
		handler, exists := ReplyHandlerMapper[reply.Type]
		if !exists {
			fmt.Println("Unknown reply:", reply.Type)
			s.OperationComplete.L.Lock()
			s.OperationComplete.Signal()
			s.OperationComplete.L.Unlock()
			continue
		}
		handler(s, reply)
	}
}

func (s *Session) Writer(wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from ", r)
		}
		s.OperationComplete.L.Lock()
		s.OperationComplete.Signal()
		s.OperationComplete.L.Unlock()
		fmt.Println("Writer Exit")
		wg.Done()
		return
	}()
	for command := range s.CommandQueue {
		handler, err := s.Context.GetHandler(command)
		if err != nil {
			fmt.Println("Unknown command:", command)
			s.OperationComplete.L.Lock()
			s.OperationComplete.Signal()
			s.OperationComplete.L.Unlock()
			continue
		}
		handler(s, command)
		s.OperationComplete.L.Lock()
		s.OperationComplete.Signal()
		s.OperationComplete.L.Unlock()
	}
}
