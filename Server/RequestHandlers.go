package main

import "fmt"

func LoginRequestHandler(session *Session, request Request) {
	//todo login logic
	fmt.Println("Login Requested")
	session.Context = NewAuthenticatedContext()
}

func RegisterRequestHandler(session *Session, request Request) {
	//todo register logic
	fmt.Println("Register Requested")
	session.Context = NewAuthenticatedContext()
}

func LogoutRequestHandler(session *Session, request Request) {
	//todo logout logic
	fmt.Println("Logout Requested")
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
