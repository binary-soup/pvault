package cmd

import (
	cmdworkflow "passwords/cmd/workflow"

	"github.com/binary-soup/go-command/command"
)

type WithdrawCommand struct {
	command.CommandBase
}

func NewWithdrawCommand() WithdrawCommand {
	return WithdrawCommand{
		CommandBase: command.NewCommandBase("withdraw", "decrypt and remove a file from the vault"),
	}
}

func (cmd WithdrawCommand) Run(args []string) error {
	return cmdworkflow.NewCmdWorkflow(cmd.Flags, args).RunDecrypt("Withdraw", true)
}
