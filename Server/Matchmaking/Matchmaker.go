package Matchmaking

import (
	"Server/Communication"
	"Server/Database"
	"slices"
)

type MatchRequestParams struct {
	MatchPlayerCount        int
	MatchPairingPreferences []string
}

func (gm *GameManager) AddPlayer(player *Database.PlayerDB, replyChannel chan Communication.Reply, mrp MatchRequestParams) {
	mp := &MatchPlayer{Player: player, ReplyChannel: replyChannel}
	gm.MatchingMutex.Lock()
	success := false
	successInd := -1
	for i := 0; i < len(gm.WaitingMatches); i++ {
		if gm.WaitingMatches[i].Capacity == mrp.MatchPlayerCount && len(gm.WaitingMatches[i].Players) < gm.WaitingMatches[i].Capacity {
			var matchDegree float32 = 0
			for j := 0; j < len(gm.WaitingMatches[i].Players); j++ {
				if slices.Contains(mrp.MatchPairingPreferences, gm.WaitingMatches[i].Players[j].Player.Username) {
					matchDegree += 1
				}
			}
			matchDegree /= float32(len(mrp.MatchPairingPreferences))
			if matchDegree > 0.75 {
				success = true
				successInd = i
				break
			}
		}
	}
	if !success {
		arr := make([]*MatchPlayer, 0, 20)
		arr = append(arr, mp)
		gm.WaitingMatches = append(gm.WaitingMatches, &Match{mrp.MatchPlayerCount, arr})
		gm.MatchingMutex.Unlock()
		return
	}
	gm.WaitingMatches[successInd].Players = append(gm.WaitingMatches[successInd].Players, mp)
	if len(gm.WaitingMatches[successInd].Players) == gm.WaitingMatches[successInd].Capacity {
		gm.MatchingMutex.Unlock()
		go gm.RunGameServer(gm.WaitingMatches[successInd])
	} else {
		gm.MatchingMutex.Unlock()
		return
	}
}

func (gm *GameManager) RemovePlayer(player *Database.PlayerDB) bool {
	gm.MatchingMutex.Lock()
	for i := 0; i < len(gm.WaitingMatches); i++ {
		ind := 0
		found := false
		for ; ind < len(gm.WaitingMatches); ind++ {
			if gm.WaitingMatches[i].Players[ind].Player.Player_id == player.Player_id {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		gm.WaitingMatches[i].Players = slices.Delete(gm.WaitingMatches[i].Players, ind, ind+1)
		if len(gm.WaitingMatches[i].Players) == 0 {
			gm.WaitingMatches = slices.Delete(gm.WaitingMatches, i, i+1)
		}
		gm.MatchingMutex.Unlock()
		return true
	}
	gm.MatchingMutex.Unlock()
	return false
}
