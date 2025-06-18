package cmd

import (
	"fmt"
	"os"
	"passwords/data"
	"passwords/data/vault"
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
		return util.Error("'name' must not be empty")
	}

	password, err := data.LoadPasswordFile(*name + ".json")
	if err != nil {
		return err
	}

	index, err := vault.LoadIndexFile(*name + ".index.json")
	if err != nil {
		index = cfg.Vault.NewIndex()
	}

	err = workflows.EncryptToVault(cfg.Vault, password, index)
	if err != nil {
		return err
	}

	if !*keep {
		os.Remove(*name + ".json")
		os.Remove(*name + ".index.json")
	}

	fmt.Printf("%s -> %s\n", ITEM_STYLE.FormatF("\"%s\"", password.Name), style.BoldInfo.Format("Saved to Vault"))
	return nil
}
