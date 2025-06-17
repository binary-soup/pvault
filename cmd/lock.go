package cmd

import (
	"fmt"
	"os"
	"passwords/data"
	"passwords/workflows"

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
	keep := cmd.Flags.Bool("keep", false, "keep the original password file")
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

	err = workflows.EncryptToVault(cfg.Vault, password, filename)
	if err != nil {
		return err
	}

	if !*keep {
		err = os.Remove(filename)
		if err != nil {
			return util.ChainError(err, "error deleting password file")
		}
	}

	fmt.Printf("%s -> %s\n", ITEM_STYLE.Format(filename), style.BoldInfo.Format("Saved to Vault"))
	return nil
}
