package Matchmaking

import (
	"Server/Communication"
	"Server/Database"
	"fmt"
	godsPq "github.com/emirpasic/gods/queues/priorityqueue"
	"slices"
	"sync"
)

type MatchRequestParams struct {
	MatchPlayerCount        int
	MatchPairingPreferences []string //todo
}

type MatchingMatches struct {
	matchScale float32
	match      *Match
}

func NewMatchingMatches(matchScale float32, match *Match) *MatchingMatches {
	return &MatchingMatches{match: match, matchScale: matchScale}
}

func MatchingMatchesComparator(a interface{}, b interface{}) int {
	aObj := a.(MatchingMatches)
	bObj := b.(MatchingMatches)
	if aObj.matchScale >= bObj.matchScale {
		return 1
	} else {
		return -1
	}
}

func (gm *GameManager) AddPlayer(player *Database.PlayerDB, replyChannel *chan Communication.Reply, replyMutex *sync.Mutex, mrp MatchRequestParams) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	fmt.Println("Matchmaking start")
	mp := &MatchPlayer{Player: player, ReplyChannel: replyChannel, ReplyMutex: replyMutex}
	gm.MatchingMutex.Lock()
	bestMatching := godsPq.NewWith(MatchingMatchesComparator)

	for i := 0; i < len(gm.WaitingMatches); i++ {
		if gm.WaitingMatches[i].Capacity == mrp.MatchPlayerCount && len(gm.WaitingMatches[i].Players) < gm.WaitingMatches[i].Capacity {
			var matchDegree float32 = 1
			for j := 0; j < len(gm.WaitingMatches[i].Players); j++ {
				if slices.Contains(mrp.MatchPairingPreferences, gm.WaitingMatches[i].Players[j].Player.Username) {
					matchDegree += 1
				}
			}
			matchDegree /= float32(len(mrp.MatchPairingPreferences) + 1)
			bestMatching.Enqueue(NewMatchingMatches(matchDegree, gm.WaitingMatches[i]))

		}
	}

	if bestMatching.Empty() {
		fmt.Println("Matchmaker didn't find a match, starts a new one")
		arr := make([]*MatchPlayer, 0, 20) //todo max player count
		arr = append(arr, mp)
		nm := &Match{mrp.MatchPlayerCount, arr}
		gm.WaitingMatches = append(gm.WaitingMatches, nm)
		gm.MatchingMutex.Unlock()
		if nm.Capacity == len(nm.Players) {
			go gm.RunGameServer(nm)
		} else {
			fmt.Println("Not starting yet")
		}
		return
	}
	rawBestFit, _ := bestMatching.Dequeue()
	bestFit := rawBestFit.(MatchingMatches).match
	fmt.Println("Match found, added to a waiting match")
	bestFit.Players = append(bestFit.Players, mp)
	if len(bestFit.Players) == bestFit.Capacity {
		gm.MatchingMutex.Unlock()
		go gm.RunGameServer(bestFit)
	} else {
		gm.MatchingMutex.Unlock()
		return
	}
}

func (gm *GameManager) RemovePlayer(player *Database.PlayerDB) bool {
	fmt.Println("matchmaking stop")
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
