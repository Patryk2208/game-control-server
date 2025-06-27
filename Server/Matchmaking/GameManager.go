package Matchmaking

import (
	"Server/Database"
	"github.com/gorilla/websocket"
	"sync"
)

type GameInstance struct {
	Game     Database.GameDB
	GameInfo Match
	//todo connections
	ControlConnection *websocket.Conn
}

type Match struct {
	Capacity int
	Players  []*Database.PlayerDB
}

type GameManager struct {
	Mutex          *sync.Mutex
	WaitingMatches []Match
	ActiveGames    []GameInstance
}

func NewGameManager() *GameManager {
	const maxMatchCount = 10000
	return &GameManager{
		WaitingMatches: make([]Match, 0, maxMatchCount),
		ActiveGames:    make([]GameInstance, 0, maxMatchCount),
	}
}
