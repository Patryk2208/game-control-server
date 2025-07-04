package GameplayManager

import (
	"github.com/google/uuid"
	"net"
	"os"
	"os/exec"
	"strconv"
)

//todo after receiving connection params start a game client and block until ended

type GameInstanceParams struct {
	GameExecutable    string
	ServerIp          net.IP
	ServerPort        int
	DedicatedPlayerId uuid.UUID
}

func StartGameplay(gip GameInstanceParams) (bool, error) {
	GameProgram := exec.Command(gip.GameExecutable,
		gip.ServerIp.String(),
		strconv.Itoa(gip.ServerPort),
		gip.DedicatedPlayerId.String())
	GameProgram.Stdin = os.Stdin
	GameProgram.Stdout = os.Stdout
	GameProgram.Stderr = os.Stderr
	err := GameProgram.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}
