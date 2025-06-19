package cmd

import (
	"fmt"
	"passwords/crypt"
	"passwords/tools"

	"github.com/binary-soup/go-command/command"
	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type GenPasswordCommand struct {
	command.CommandBase
}

func NewGenPasswordCommand() GenPasswordCommand {
	return GenPasswordCommand{
		CommandBase: command.NewCommandBase("gen", "generate a new random password"),
	}
}

func (cmd GenPasswordCommand) Run(args []string) error {
	len := cmd.Flags.Int("l", 0, "length of the password")
	cmd.Flags.Parse(args)

	if *len < 1 {
		return util.Error("(l)ength missing or invalid")
	}

	password, err := crypt.RandPassword(*len)
	if err != nil {
		return util.ChainError(err, "error generating password")
	}

	err = tools.CopyToClipboard(password)
	if err != nil {
		return err
	}

	fmt.Printf("%s -> %s\n", style.Bolded.FormatF("PASSWORD[%d]", *len), style.BoldInfo.Format("Copied to Clipboard"))
	return nil
}
