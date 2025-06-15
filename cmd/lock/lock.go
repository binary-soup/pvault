package cmdlock

import (
	"os"
	"passwords/data"
	"path/filepath"

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

	cfg, err := data.LoadConfig()
	if err != nil {
		return err
	}

	if *path == "" {
		return util.Error("path must not be empty")
	}

	password, err := data.LoadPasswordFile(*path)
	if err != nil {
		return err
	}

	err = cfg.Vault.SavePassword(password, filepath.Base(*path))
	if err != nil {
		return err
	}

	err = os.Remove(*path)
	if err != nil {
		return util.ChainError(err, "error deleting password file")
	}

	style.New(style.Yellow).Print(*path)
	style.BoldInfo.Println(" -> Saved to Vault.")

	return nil
}
