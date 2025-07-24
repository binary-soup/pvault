package cmd

import (
	"fmt"
	"pvault/data/version"

	"github.com/binary-soup/go-commando/command"
	"github.com/binary-soup/go-commando/style"
)

type VersionCommand struct {
	command.CommandBase
}

func NewVersionCommand() VersionCommand {
	return VersionCommand{
		CommandBase: command.NewCommandBase("version", "report version information"),
	}
}

func (cmd VersionCommand) Run(args []string) error {
	fmt.Printf("App:\t%s\n", style.Bolded.Format(version.APP))
	fmt.Printf("Vault:\t%s\n", style.Bolded.FormatF("%d", version.VAULT))
	return nil
}
