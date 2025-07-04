package main

type Command struct {
	Name   string
	Args   []string
	Source string
}

type CommandHandler func(session *Session, Command Command)
