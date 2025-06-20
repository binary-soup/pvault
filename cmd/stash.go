package cmd

import (
	cmdworkflow "passwords/cmd/workflow"

	"github.com/binary-soup/go-command/command"
)

type StashCommand struct {
	command.CommandBase
}

func NewStashCommand() StashCommand {
	return StashCommand{
		CommandBase: command.NewCommandBase("stash", "encrypt and stash a new file in the vault"),
	}
}

func (cmd StashCommand) Run(args []string) error {
	return cmdworkflow.NewCmdWorkflow(cmd.Flags, args).RunEncrypt("Stash", true)
}
