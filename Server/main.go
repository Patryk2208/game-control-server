package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			panic(err)
		}
	}(c)

	session := NewSession()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go session.Reader()
	go session.Writer()
	wg.Wait()

	fmt.Println("Client connection closed")
}

func basicListen() {
	http.HandleFunc("/ws", handleConnection)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func main() {
	basicListen()
}
