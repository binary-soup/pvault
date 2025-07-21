package cmd

import (
	"flag"
	"fmt"
	"pvault/crypt"
	"pvault/tools"

	"github.com/binary-soup/go-commando/alert"
	"github.com/binary-soup/go-commando/command"
	"github.com/binary-soup/go-commando/style"
)

type GenPasswordCommand struct {
	command.CommandBase
	flags *genPasswordFlags
}

type genPasswordFlags struct {
	Length *int
}

func (f *genPasswordFlags) Set(flags *flag.FlagSet) {
	f.Length = flags.Int("l", 0, "length of the password")
}

func NewGenPasswordCommand() GenPasswordCommand {
	flags := new(genPasswordFlags)

	return GenPasswordCommand{
		CommandBase: command.NewCommandBase("gen", "generate a new random password", flags),
		flags:       flags,
	}
}

func (cmd GenPasswordCommand) Run() error {
	if *cmd.flags.Length < 1 {
		return alert.Error("(l)ength missing or invalid")
	}

	password, err := crypt.RandPassword(*cmd.flags.Length)
	if err != nil {
		return alert.ChainError(err, "error generating password")
	}

	err = tools.CopyToClipboard(password)
	if err != nil {
		return err
	}

	fmt.Printf("%s -> %s\n", style.Bolded.FormatF("PASSWORD[%d]", *cmd.flags.Length), style.BoldInfo.Format("Copied to Clipboard"))
	return nil
}
