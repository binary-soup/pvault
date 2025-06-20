package cmd

import (
	cmdworkflow "pvault/cmd/workflow"

	"github.com/binary-soup/go-command/command"
)

type UnlockCommand struct {
	command.CommandBase
}

func NewUnlockCommand() UnlockCommand {
	return UnlockCommand{
		CommandBase: command.NewCommandBase("unlock", "temporarily decrypt a file from the vault"),
	}
}

func (cmd UnlockCommand) Run(args []string) error {
	return cmdworkflow.NewCmdWorkflow(cmd.Flags, args).RunDecrypt("Unlock", false)
}
