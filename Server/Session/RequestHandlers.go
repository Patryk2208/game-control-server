package Session

import (
	"Server/Communication"
	"Server/Matchmaking"
	"fmt"
	"strconv"
	"strings"
)

func LoginRequestHandler(session *Session, request Communication.Request) {
	data := strings.Split(request.Message, " ")
	if len(data) != 3 {
		session.ReplyQueue <- Communication.Reply{Type: Communication.SystemReply, Message: "F"}
		return
	}
	success, player, err := session.DbPool.TryLogin(data[1], data[2])
	if err != nil || !success {
		session.ReplyQueue <- Communication.Reply{Type: Communication.SystemReply, Message: "F"}
		return
	}
	session.Player = player
	session.ReplyQueue <- Communication.Reply{Type: Communication.SystemReply, Message: "T"}
	session.Context = NewAuthenticatedContext()
}

func RegisterRequestHandler(session *Session, request Communication.Request) {
	data := strings.Split(request.Message, " ")
	if len(data) != 3 {
		session.ReplyQueue <- Communication.Reply{Type: Communication.SystemReply, Message: "F"}
		return
	}
	success, err := session.DbPool.TryRegisterUser(data[1], data[2])
	if err != nil || !success {
		session.ReplyQueue <- Communication.Reply{Type: Communication.SystemReply, Message: "F"}
		return
	}
	session.ReplyQueue <- Communication.Reply{Type: Communication.SystemReply, Message: "T"}
}

func LogoutRequestHandler(session *Session, request Communication.Request) {
	success, err := session.DbPool.TryLogout(session.Player)
	if err != nil || !success {
		session.ReplyQueue <- Communication.Reply{Type: Communication.SystemReply, Message: "F"}
		return
	}
	session.ReplyQueue <- Communication.Reply{Type: Communication.SystemReply, Message: "T"}
	session.Context = NewNormalContext()
}

func StartGameRequestHandler(session *Session, request Communication.Request) {
	mrp := CreateMatchRequestParams(request)
	session.GameManager.AddPlayer(session.Player, session.ReplyQueue, mrp)
	session.ReplyQueue <- Communication.Reply{Type: Communication.SystemReply, Message: "T"}
	session.Context = NewWaitingContext()
}

func CreateMatchRequestParams(request Communication.Request) Matchmaking.MatchRequestParams {
	data := strings.Split(request.Message, " ")
	if len(data) <= 1 {
		return Matchmaking.MatchRequestParams{MatchPlayerCount: 1, MatchPairingPreferences: nil}
	}
	pairing := make([]string, 0, 20)
	n, err := strconv.Atoi(data[1])
	if err != nil {
		for i := 2; i < len(data); i++ {
			pairing = append(pairing, data[i])
		}
		return Matchmaking.MatchRequestParams{MatchPlayerCount: len(pairing) + 1, MatchPairingPreferences: pairing}
	}
	for i := 3; i < len(data); i++ {
		pairing = append(pairing, data[i])
	}
	return Matchmaking.MatchRequestParams{MatchPlayerCount: n, MatchPairingPreferences: pairing}
}

func EndGameRequestHandler(session *Session, request Communication.Request) {
	//todo end game logic
	fmt.Println("Game End Requested")
	session.Context = NewAuthenticatedContext()
}

func ExitRequestHandler(session *Session, request Communication.Request) {
	fmt.Println("Exit Requested")
	panic(fmt.Errorf("exit requested"))
}

func ExitWithLogoutRequestHandler(session *Session, request Communication.Request) {
	LogoutRequestHandler(session, request)
	ExitRequestHandler(session, request)
}

func StopWaitingRequestHandler(session *Session, request Communication.Request) {
	success := session.GameManager.RemovePlayer(session.Player)
	if success {
		session.ReplyQueue <- Communication.Reply{Type: Communication.SystemReply, Message: "T"}
		session.Context = NewAuthenticatedContext()
	} else {
		session.ReplyQueue <- Communication.Reply{Type: Communication.SystemReply, Message: "F"}
	}
}

func StopWaitingAndLogoutRequestHandler(session *Session, request Communication.Request) {
	StopWaitingRequestHandler(session, request)
	LogoutRequestHandler(session, request)
}

func StopWaitingAndExitRequestHandler(session *Session, request Communication.Request) {
	StopWaitingRequestHandler(session, request)
	ExitRequestHandler(session, request)
}
