package cmd

import (
	cmdworkflow "pvault/cmd/workflow"

	"github.com/binary-soup/go-command/command"
)

type RelockCommand struct {
	command.CommandBase
}

func NewRelockCommand() RelockCommand {
	return RelockCommand{
		CommandBase: command.NewCommandBase("relock", "re-encrypt a file back into the vault"),
	}
}

func (cmd RelockCommand) Run(args []string) error {
	return cmdworkflow.NewCmdWorkflow(cmd.Flags, args).RunEncrypt("Relock", false)
}
