package main

import (
	"fmt"
)

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
}

func ParseRegisterReply(reply Reply) bool {
	if reply.Message[0] == 'T' {
		fmt.Println("Registration successful")
		return true
	} else {
		fmt.Println("Registration failed")
		return false
	}
}

func StartGameCommandHandler(session *Session, command Command) {
	message := command.Name
	for i := 0; i < len(command.Args); i++ {
		message += " " + command.Args[i]
	}
	err := session.ServerConn.WriteMessage(1, []byte(message))
	if err != nil {
		return
	}
	reply := <-session.ReplyChannel
	if !ParsePlayReply(reply) {
		return
	}
	session.Context = NewWaitingContext()
}

func ParsePlayReply(reply Reply) bool {
	if reply.Message[0] == 'T' {
		fmt.Println("Waiting for a game to begin")
		return true
	} else {
		fmt.Println("Play command failed")
		return false
	}
}

func StopWaitingCommandHandler(session *Session, command Command) {
	err := session.ServerConn.WriteMessage(1, []byte(command.Name))
	if err != nil {
		return
	}
	reply := <-session.ReplyChannel
	if !ParseStopWaitingReply(reply) {
		return
	}
	session.Context = NewAuthenticatedContext()
}

func ParseStopWaitingReply(reply Reply) bool {
	if reply.Message[0] == 'T' {
		fmt.Println("Stopped waiting")
		return true
	} else {
		fmt.Println("still waiting")
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
