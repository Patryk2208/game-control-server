package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (s *Session) StartREPL() {
	reader := bufio.NewReader(os.Stdin)
	for {
		s.DisplayPrompt()
		input, _ := reader.ReadString('\n')
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
