package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

func main() {
	session := NewSession()
	defer session.Close()

	c, _, err := websocket.DefaultDialer.Dial("ws://10.111.231.112:8080/ws", nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}

	session.ServerConn = c
	wg := sync.WaitGroup{}
	wg.Add(2)
	go session.Reader(&wg)
	go session.Writer(&wg)

	fmt.Println("Game CLI Session Started (v1.0)")
	fmt.Println("Type 'help' for available commands")
	session.StartREPL()

	wg.Wait()
}
