package main

import (
	"flag"
	"os"
	"pvault/cmd"

	"github.com/binary-soup/go-command/alert"
	"github.com/binary-soup/go-command/command"
)

func main() {
	ls := flag.Bool("ls", false, "list all commands")
	flag.Parse()

	runner := command.NewRunner(
		cmd.NewSearchCommand(),
		cmd.NewStashCommand(),
		cmd.NewWithdrawCommand(),
		cmd.NewDeleteCommand(),
		cmd.NewUnlockCommand(),
		cmd.NewRelockCommand(),
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
		alert.Print(err)
	}
}
