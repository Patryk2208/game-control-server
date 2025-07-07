package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (s *Session) Start() {
	reader := bufio.NewReader(os.Stdin)
	for {
		s.DisplayPrompt()

		s.CommandQueue <- cmd
		s.OperationComplete.L.Lock()
		s.OperationComplete.Wait()
		s.OperationComplete.L.Unlock()
		if cmd.Name == "exit" {
			break
		}
	}
}

func (s *Session) DisplayPrompt() {
	fmt.Printf("[" + s.Context.GetPrompt() + "]> ")
}

func Console(reader *bufio.Reader) bool {
	input, err := reader.ReadString('\n')
	if err != nil {
		return true
	}
	input = strings.TrimSpace(input)

	if input == "" {
		return true
	}

	parts := strings.Fields(input)
	cmd := Command{
		Name:   parts[0],
		Args:   parts[1:],
		Source: "user",
	}
}

func Game() {
	
}
