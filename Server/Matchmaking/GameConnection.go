package Matchmaking

import (
	"Server/Communication"
	"Server/Database"
	"context"
	"fmt"
	"slices"
)

func (gm *GameManager) RunGameServer(readyGame *Match) {
	gm.MatchingMutex.Lock()
	ind := slices.Index(gm.WaitingMatches, readyGame)
	if ind == -1 || len(readyGame.Players) != readyGame.Capacity {
		gm.MatchingMutex.Unlock()
		return
	}
	gm.WaitingMatches = slices.Delete(gm.WaitingMatches, ind, ind+1)
	gm.MatchingMutex.Unlock()

	ctx := context.Background()
	ip, port, err := gm.CreateGameServer(ctx)
	if err != nil {
		panic(err)
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
		player.ReplyChannel <- universalReply
	}

	gm.WatchContainerState(ctx, gm.ContainerInfo.StdGSTemplate)
	gm.ActiveMutex.Lock()
	//todo move to archived in db
	ind = slices.Index(gm.ActiveGames, gi)
	if ind != -1 {
		panic("game already closed")
	}
	gm.ActiveGames = slices.Delete(gm.ActiveGames, ind, ind+1)
	gm.ActiveMutex.Unlock()
}
