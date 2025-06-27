package Session

import (
	"Server/Database"
	"Server/Matchmaking"
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
	"sync"
)

type Session struct {
	Context     UserConnectionContext
	Player      *Database.PlayerDB
	ClientConn  *websocket.Conn
	ReplyQueue  chan Reply
	DbPool      *Database.DBConnectionPool
	GameManager *Matchmaking.GameManager
}

func NewSession(c *websocket.Conn, pool *Database.DBConnectionPool, gm *Matchmaking.GameManager) *Session {
	return &Session{
		Context:     NewNormalContext(),
		ClientConn:  c,
		Player:      nil,
		ReplyQueue:  make(chan Reply),
		DbPool:      pool,
		GameManager: gm,
	}
}

func (s *Session) Reader(wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from ", r)
		}
		s.ReplyQueue <- Reply{Type: ExitReply, Message: "Exit"}
		fmt.Println("Reader Exit")
		wg.Done()
		return
	}()
	for {
		msgType, msg, err := s.ClientConn.ReadMessage()
		if err != nil {
			fmt.Println("Reader Exit")
			panic(err)
		}
		temp := strings.Split(string(msg), " ")[0]
		requestType, ok := RequestTypeMapper[temp]
		if !ok {
			continue
		}
		req := NewRequest(requestType, string(msg))
		handler, err := s.Context.GetHandler(req)
		if err != nil {
			continue
		}
		handler(s, *req)
		fmt.Printf("Message type: %d, %s\n", msgType, string(msg))
	}
}

func (s *Session) Writer(wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from ", r)
		}
		err := s.ClientConn.WriteMessage(1, []byte("T exit OK"))
		if err != nil {
		}
		fmt.Println("Writer Exit")
		wg.Done()
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
