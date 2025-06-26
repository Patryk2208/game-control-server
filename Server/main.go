package main

import (
	"Server/Database"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"math"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var globalConnectionPool *Database.DBConnectionPool = nil

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

	session := NewSession(c, globalConnectionPool)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go session.Reader(&wg)
	go session.Writer(&wg)
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
	config := Database.ConnectionConfig{
		Host:            "172.17.0.2",
		Port:            5432,
		Database:        "users",
		Username:        "patryk",
		Password:        "sql",
		MaxConnections:  int32(4 * runtime.NumCPU()),
		MinConnections:  int32(math.Round(0.2 * 4 * float64(runtime.NumCPU()))),
		MaxConnIdleTime: 5 * time.Minute,
		MaxConnLifetime: 30 * time.Minute,
	}
	var err error
	globalConnectionPool, err = Database.InitConnectionPool(context.Background(), config)
	if err != nil {
		panic(err)
	}
	defer Database.CloseConnectionPool(globalConnectionPool)
	basicListen()
}
