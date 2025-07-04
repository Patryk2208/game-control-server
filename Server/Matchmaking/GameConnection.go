package Matchmaking

import "slices"

func (gm *GameManager) StartGame(readyGame *Match) {
	gm.MatchingMutex.Lock()
	ind := slices.Index(gm.WaitingMatches, readyGame)
	if ind == -1 || len(readyGame.Players) != readyGame.Capacity {
		gm.MatchingMutex.Unlock()
		return
	}
	gm.WaitingMatches = slices.Delete(gm.WaitingMatches, ind, ind+1)
	gm.MatchingMutex.Unlock()
	gm.ActiveMutex.Lock()
	//todo run a game container add it to active games
	gm.ActiveMutex.Unlock()

	//todo maintain the control connection, send its ip to all players, maintain container

	gm.ActiveMutex.Lock()
	//todo remove from active connections
	gm.ActiveMutex.Unlock()
}
