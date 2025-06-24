package main

type Command struct {
	Name   string
	Args   []string
	Source string
}

type CommandHandler func(session *Session, Command Command)

/*func handleExit(s *Session, args []string) error {
	s.Close()
	os.Exit(0)
	return nil
}

func handleHelp(s *Session, args []string) error {
	//TODO help
	return nil
}

func handleStatus(s *Session, args []string) error {
	s.ReplyChannel <- Reply{SystemMessage, s.State.String()}
	return nil
}
*/
