package main

import (
	"fmt"
)

//todo

func LoginCommandHandler(session *Session, command Command) {
	message := command.Name
	for i := 0; i < len(command.Args); i++ {
		message += " " + command.Args[i]
	}
	err := session.ServerConn.WriteMessage(1, []byte(message))
	if err != nil {
		return
	}
	reply := <-session.ReplyChannel
	if !ParseLoginReply(reply) {
		return
	}
	session.Context = NewAuthenticatedContext()
}

func ParseLoginReply(reply Reply) bool {
	if reply.Message[0] == 'T' {
		fmt.Println("Authentication successful")
		return true
	} else {
		fmt.Println("Authentication failed")
		return false
	}
}

func RegisterCommandHandler(session *Session, command Command) {
	fmt.Println("register command")
}

func NormalHelpCommandHandler(session *Session, command Command) {
	fmt.Println("help command")
}

func LogoutCommandHandler(session *Session, command Command) {
	fmt.Println("logout command")
}

func AuthenticatedHelpCommandHandler(session *Session, command Command) {
	fmt.Println("help command")
}

func ExitCommandHandler(session *Session, command Command) {
	err := session.ServerConn.WriteMessage(1, []byte(command.Name))
	if err != nil {
		panic(fmt.Errorf("exit failure"))
	}
	reply := <-session.ReplyChannel
	if reply.Message[0] == 'T' {
		panic(fmt.Errorf("exit success"))
	} else {
		panic(fmt.Errorf("exit failure"))
	}
}

func PlayCommandHandler(session *Session, command Command) {
	fmt.Println("play command")
}
