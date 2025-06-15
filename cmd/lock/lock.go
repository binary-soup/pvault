package cmdlock

import (
	"passwords/data"

	"github.com/binary-soup/go-command/command"
	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type LockCommand struct {
	command.CommandBase
}

func NewLockCommand() LockCommand {
	return LockCommand{
		CommandBase: command.NewCommandBase("lock", "encrypt the provided password file"),
	}
}

func (cmd LockCommand) Run(args []string) error {
	path := cmd.Flags.String("path", "", "path to the password file")
	cmd.Flags.Parse(args)

	if *path == "" {
		return util.Error("path must not be empty")
	}

	password, err := data.LoadPasswordFile(*path)
	if err != nil {
		return err
	}

	output := *path + ".crypt"

	err = password.EncryptToFile(output)
	if err != nil {
		return err
	}

	style.Create.PrintF("+ %s\n", output)
	return nil
}
