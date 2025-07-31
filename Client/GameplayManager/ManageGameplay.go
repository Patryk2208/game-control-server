package GameplayManager

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
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
		GameExecutable: "",
		ServerIp:       ip,
		ServerPort:     port,
	}
	if runtime.GOOS == "linux" {
		gip.GameExecutable = "./Client"
	} else if runtime.GOOS == "windows" {
		gip.GameExecutable = ".\\Client.exe"
	} else {
		return
	}

	StartGameplay(gip)
}

func ManageInput(cmd *exec.Cmd) {
	stdinPipe, _ := cmd.StdinPipe()
	inputChan := make(chan byte, 32)

	if runtime.GOOS == "windows" {
		exec.Command("cmd", "/c", "REG ADD HKCU\\CONSOLE /v QuickEdit /t REG_DWORD /d 0 /f").Run()
	}

	go func() {
		buf := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				close(inputChan)
				return
			}
			inputChan <- buf[0]
		}
	}()

	go func() {
		for b := range inputChan {
			stdinPipe.Write([]byte{b})
		}
	}()
}

func StartGameplay(gip GameInstanceParams) (bool, error) {
	GameProgram := exec.Command(gip.GameExecutable,
		gip.ServerIp,
		strconv.Itoa(gip.ServerPort))
	GameProgram.Stdout = os.Stdout
	GameProgram.Stderr = os.Stderr

	ManageInput(GameProgram)
	err := GameProgram.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}
