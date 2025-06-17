package cmd

import (
	"fmt"
	"passwords/data"
	"passwords/workflows"

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
	name := cmd.Flags.String("name", "", "name of the vault item")
	cmd.Flags.Parse(args)

	cfg, err := data.LoadConfig()
	if err != nil {
		return err
	}

	if *name == "" {
		return util.Error("name must not be empty")
	}
	filename := *name + ".json"

	password, err := workflows.DecryptFromVault(cfg.Vault, filename)
	if err != nil {
		return err
	}

	err = password.SaveToFile(filename)
	if err != nil {
		return err
	}

	err = workflows.DeleteFromVault(cfg.Vault, filename)
	if err != nil {
		return err
	}

	fmt.Printf("%s -> %s\n", ITEM_STYLE.Format(filename), style.BoldInfo.Format("Loaded to Vault"))
	return nil
}
