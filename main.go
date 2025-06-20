package main

import (
	"flag"
	"fmt"
	"os"
	"passwords/cmd"

	"github.com/binary-soup/go-command/command"
	"github.com/binary-soup/go-command/style"
)

func main() {
	ls := flag.Bool("ls", false, "list all commands")
	flag.Parse()

	runner := command.NewRunner(
		cmd.NewSearchCommand(),
		cmd.NewStashCommand(),
		cmd.NewWithdrawCommand(),
		cmd.NewCopyCommand(),
		cmd.NewImportCommand(),
		cmd.NewGenPasswordCommand(),
		cmd.NewSyncCommand(),
	)

	if *ls || len(os.Args) < 2 {
		runner.ListCommands()
		return
	}

	if err := runner.RunCommand(os.Args[1], os.Args[2:]); err != nil {
		style.BoldError.Print("ERROR: ")
		fmt.Println(err)
	}
}
