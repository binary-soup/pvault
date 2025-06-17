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
	id := cmd.Flags.Uint("id", 0, "id of the vault item (use 'search' to query by name)")
	out := cmd.Flags.String("out", "", "name of the out file")
	cmd.Flags.Parse(args)

	cfg, err := data.LoadConfig()
	if err != nil {
		return err
	}

	if *id == 0 {
		return util.Error("'id' missing or invalid")
	}

	if *out == "" {
		return util.Error("'out' must not be empty")
	}

	password, index, err := workflows.DecryptFromVault(cfg.Vault, *id)
	if err != nil {
		return err
	}

	err = password.SaveToFile(*out + ".json")
	if err != nil {
		return err
	}

	err = index.SaveToFile(*out + ".index.json")
	if err != nil {
		return err
	}

	err = workflows.DeleteFromVault(cfg.Vault, *id)
	if err != nil {
		return err
	}

	fmt.Printf("%s -> %s\n", ITEM_STYLE.FormatF("\"%s\"", password.Name), style.BoldInfo.Format("Loaded from Vault"))
	return nil
}
