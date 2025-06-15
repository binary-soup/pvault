package cmdunlock

import (
	"passwords/data"
	"strings"

	"github.com/binary-soup/go-command/command"
	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type UnlockCommand struct {
	command.CommandBase
}

func NewUnlockCommand() UnlockCommand {
	return UnlockCommand{
		CommandBase: command.NewCommandBase("unlock", "decrypt the provided crypt file"),
	}
}

func (cmd UnlockCommand) Run(args []string) error {
	path := cmd.Flags.String("path", "", "path to the crypt file")
	cmd.Flags.Parse(args)

	if *path == "" {
		return util.Error("path must not be empty")
	}

	password, err := data.DecryptPasswordFromFile(*path)
	if err != nil {
		return err
	}

	output, _ := strings.CutSuffix(*path, ".crypt")

	err = password.SaveToFile(output)
	if err != nil {
		return err
	}

	style.Create.PrintF("+ %s\n", output)
	return nil
}
