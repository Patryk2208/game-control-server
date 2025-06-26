package main

import (
	"fmt"
	"strings"
)

func LoginRequestHandler(session *Session, request Request) {
	data := strings.Split(request.Message, " ")
	if len(data) != 3 {
		session.ReplyQueue <- Reply{SystemReply, "F"}
		return
	}
	success, player, err := session.DbPool.TryLogin(data[1], data[2])
	if err != nil || !success {
		session.ReplyQueue <- Reply{SystemReply, "F"}
		return
	}
	session.Player = player
	session.ReplyQueue <- Reply{SystemReply, "T"}
	session.Context = NewAuthenticatedContext()
}

func RegisterRequestHandler(session *Session, request Request) {
	data := strings.Split(request.Message, " ")
	if len(data) != 3 {
		session.ReplyQueue <- Reply{SystemReply, "F"}
		return
	}
	success, err := session.DbPool.TryRegisterUser(data[1], data[2])
	if err != nil || !success {
		session.ReplyQueue <- Reply{SystemReply, "F"}
		return
	}
	session.ReplyQueue <- Reply{SystemReply, "T"}
}

func LogoutRequestHandler(session *Session, request Request) {
	success, err := session.DbPool.TryLogout(session.Player)
	if err != nil || !success {
		session.ReplyQueue <- Reply{SystemReply, "F"}
		return
	}
	session.ReplyQueue <- Reply{SystemReply, "T"}
	session.Context = NewNormalContext()
}

func StartGameRequestHandler(session *Session, request Request) {
	//todo start game logic
	fmt.Println("Game Start Requested")
	session.Context = NewPlayingContext()
}

func EndGameRequestHandler(session *Session, request Request) {
	//todo end game logic
	fmt.Println("Game End Requested")
	session.Context = NewAuthenticatedContext()
}

func ExitRequestHandler(session *Session, request Request) {
	fmt.Println("Exit Requested")
	panic(fmt.Errorf("exit requested"))
}

func ExitWithLogoutRequestHandler(session *Session, request Request) {
	LogoutRequestHandler(session, request)
	ExitRequestHandler(session, request)
}
