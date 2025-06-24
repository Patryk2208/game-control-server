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
	message := command.Name
	for i := 0; i < len(command.Args); i++ {
		message += " " + command.Args[i]
	}
	err := session.ServerConn.WriteMessage(1, []byte(message))
	if err != nil {
		return
	}
	reply := <-session.ReplyChannel
	if !ParseRegisterReply(reply) {
		return
	}
	session.Context = NewAuthenticatedContext()
}

func ParseRegisterReply(reply Reply) bool {
	if reply.Message[0] == 'T' {
		fmt.Println("Authentication successful")
		return true
	} else {
		fmt.Println("Authentication failed")
		return false
	}
}

func NormalHelpCommandHandler(session *Session, command Command) {
	fmt.Println("You are not logged in")
	fmt.Println("Usage:")
	fmt.Println("help")
	fmt.Println("register <username> <password>")
	fmt.Println("login <username> <password>")
	fmt.Println("exit")
}

func LogoutCommandHandler(session *Session, command Command) {
	err := session.ServerConn.WriteMessage(1, []byte(command.Name))
	if err != nil {
		return
	}
	reply := <-session.ReplyChannel
	if !ParseLogoutReply(reply) {
		return
	}
	session.Context = NewNormalContext()
}

func ParseLogoutReply(reply Reply) bool {
	if reply.Message[0] == 'T' {
		fmt.Println("Logout successful")
		return true
	} else {
		fmt.Println("Logout failed")
		return false
	}
}

func AuthenticatedHelpCommandHandler(session *Session, command Command) {
	fmt.Println("You are logged in")
	fmt.Println("Usage:")
	fmt.Println("help")
	fmt.Println("logout")
	fmt.Println("start :: start new singleplayer game")
	fmt.Println("exit")
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
