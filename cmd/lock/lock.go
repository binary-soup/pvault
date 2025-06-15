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
	name := cmd.Flags.String("name", "", "name of the password file")
	cmd.Flags.Parse(args)

	cfg, err := data.LoadConfig()
	if err != nil {
		return err
	}

	if *name == "" {
		return util.Error("name must not be empty")
	}
	filename := *name + ".json"

	password, err := data.LoadPasswordFile(filename)
	if err != nil {
		return err
	}

	err = cfg.Vault.SavePassword(password, filepath.Base(filename))
	if err != nil {
		return err
	}

	err = os.Remove(filename)
	if err != nil {
		return util.ChainError(err, "error deleting password file")
	}

	style.New(style.Yellow).Print(filename)
	style.BoldInfo.Println(" -> Saved to Vault.")

	return nil
}
