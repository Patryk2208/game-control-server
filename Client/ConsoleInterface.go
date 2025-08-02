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
		input, err := reader.ReadString('\n') //todo
		if err != nil {
			continue
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		cmd := Command{
			Name:   parts[0],
			Args:   parts[1:],
			Source: "user",
		}

		s.OperationComplete.L.Lock()
		s.CommandQueue <- cmd
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
