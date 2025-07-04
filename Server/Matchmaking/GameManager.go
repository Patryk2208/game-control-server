package Matchmaking

import (
	"Server/Database"
	"github.com/gorilla/websocket"
	"net"
	"sync"
)

type GameContainerAddress struct {
	Ip   net.IP
	Port uint16
}

type GameInstance struct {
	Game              Database.GameDB
	GameInfo          Match
	ControlConnection *websocket.Conn
	GameAddress       GameContainerAddress
}

type Match struct {
	Capacity int
	Players  []*Database.PlayerDB
}

type GameManager struct {
	MatchingMutex  *sync.Mutex
	ActiveMutex    *sync.Mutex
	WaitingMatches []*Match
	ActiveGames    []*GameInstance
}

func NewGameManager() *GameManager {
	const maxMatchCount = 10000
	return &GameManager{
		WaitingMatches: make([]*Match, 0, maxMatchCount),
		ActiveGames:    make([]*GameInstance, 0, maxMatchCount),
	}
}
