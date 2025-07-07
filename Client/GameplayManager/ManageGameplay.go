package GameplayManager

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//todo after receiving connection params start a game client and block until ended

type GameInstanceParams struct {
	GameExecutable string
	ServerIp       string
	ServerPort     int
}

func RunGameplay(message string) {
	params := strings.Split(message, " ")
	ip := params[0]
	port, err := strconv.Atoi(params[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	gip := GameInstanceParams{
		GameExecutable: "", //todo
		ServerIp:       ip,
		ServerPort:     port,
	}
	StartGameplay(gip)
}

func StartGameplay(gip GameInstanceParams) (bool, error) {
	GameProgram := exec.Command(gip.GameExecutable,
		gip.ServerIp,
		strconv.Itoa(gip.ServerPort))
	GameProgram.Stdin = os.Stdin
	GameProgram.Stdout = os.Stdout
	GameProgram.Stderr = os.Stderr
	err := GameProgram.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}
