package cmdunlock

import (
	"os"
	"passwords/data"

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
	input := *path + ".crypt"

	password, err := data.DecryptPasswordFromFile(input)
	if err != nil {
		return err
	}

	err = password.SaveToFile(*path)
	if err != nil {
		return err
	}
	style.Create.PrintF("+ %s\n", *path)

	err = os.Remove(input)
	if err != nil {
		return util.ChainError(err, "error deleting crypt file")
	}
	style.Delete.PrintF("- %s\n", input)

	return nil
}
