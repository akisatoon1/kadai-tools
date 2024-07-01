package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	var l = log.New(os.Stderr, "error: ", 0)
	err := processCommand()
	if err != nil {
		l.Fatal(err)
	}
}

func processCommand() error {
	if len(os.Args) == 1 {
		return fmt.Errorf("引数が足りません")
	}
	switch cmd := os.Args[1]; cmd {
	case "make":
		err := make(os.Args[2:])
		if err != nil {
			return err
		}
	case "debug":
		err := debug(os.Args[2:])
		if err != nil {
			return err
		}
	case "submit":
		err := submit(os.Args[2:])
		if err != nil {
			return err
		}
	default:
		emsg := fmt.Sprintf("コマンド'%v'は存在しません\n", cmd)
		emsg += summaryCommandList([]string{"make", "debug", "submit"})
		return fmt.Errorf(emsg)
	}
	return nil
}

func summaryCommandList(list []string) string {
	s := "コマンド一覧\n\n"
	for _, cmd := range list {
		s += cmd + "\n"
	}
	return s
}
