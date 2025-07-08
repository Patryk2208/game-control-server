package Matchmaking

import (
	"Server/Communication"
	"Server/Database"
	"context"
	"fmt"
	"slices"
)

func (gm *GameManager) RunGameServer(readyGame *Match) {
	fmt.Println("tries to run game server")
	gm.MatchingMutex.Lock()
	ind := slices.Index(gm.WaitingMatches, readyGame)
	if ind == -1 || len(readyGame.Players) != readyGame.Capacity {
		fmt.Println("not enough players or game server already started")
		gm.MatchingMutex.Unlock()
		return
	}
	gm.WaitingMatches = slices.Delete(gm.WaitingMatches, ind, ind+1)
	gm.MatchingMutex.Unlock()

	ctx := context.Background()
	ip, port, err := gm.AllocateGameServer(ctx)
	fmt.Println("game server created")
	if err != nil {
		fmt.Println(err)
		return
	}
	//todo add game instance to db

	gi := &GameInstance{
		Game:              Database.GameDB{},
		GameInfo:          *readyGame,
		ControlConnection: nil,
		GameAddress: GameContainerAddress{
			Ip:   ip,
			Port: port,
		},
	}
	gm.ActiveMutex.Lock()
	gm.ActiveGames = append(gm.ActiveGames, gi)
	gm.ActiveMutex.Unlock()

	universalReply := Communication.Reply{
		Type:    Communication.GameReply,
		Message: fmt.Sprintf("G %s %d", ip, port),
	}
	for _, player := range readyGame.Players {
		player.ReplyMutex.Lock()
		player.ReplyChannel <- universalReply
		player.ReplyMutex.Unlock()
	}
	fmt.Println("game address sent to each player")

	gm.WatchContainerState(ctx)
	fmt.Println("game server ended")
	gm.ActiveMutex.Lock()
	//todo move to archived in db
	ind = slices.Index(gm.ActiveGames, gi)
	if ind != -1 {
		fmt.Println("game already closed")
		return
	}
	gm.ActiveGames = slices.Delete(gm.ActiveGames, ind, ind+1)
	gm.ActiveMutex.Unlock()
	fmt.Println("game instance completed")
}
