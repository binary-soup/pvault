package cmd

import (
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
	style.Bolded.PrintF("v%s\n", version.APP)
	return nil
}
