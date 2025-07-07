package Matchmaking

import (
	"Server/Database"
	"context"
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

	//todo maintain the control connection, send its ip to all players, maintain container
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
	//todo run a game container add it to active games
	i := len(gm.ActiveGames)
	gm.ActiveGames = append(gm.ActiveGames, gi)
	j := len(gm.ActiveGames)
	gm.ActiveMutex.Unlock()

	//todo send to client
	gm.WatchContainerState(ctx, gm.ContainerInfo.StdGSTemplate)

	gm.ActiveMutex.Lock()
	//todo move to archived in db
	gm.ActiveGames = slices.Delete(gm.ActiveGames, i, j)
	gm.ActiveMutex.Unlock()
}
