package cmd

import (
	"pvault/data/config"
	vw "pvault/workflows/vault"

	"github.com/binary-soup/go-commando/command"
)

type SearchCommand struct {
	command.ConfigCommandBase[config.Config]
}

func NewSearchCommand() SearchCommand {
	return SearchCommand{
		ConfigCommandBase: command.NewConfigCommandBase[config.Config]("search", "search the items in the vault"),
	}
}

func (cmd SearchCommand) Run(args []string) error {
	search := cmd.Flags.String("s", "", "the search term. Leave blank to list all")
	cmd.Flags.Parse(args)

	cfg, err := cmd.LoadConfig(DATA_DIR)
	if err != nil {
		return err
	}

	vw.NewVaultWorkflow(cfg.Vault).Search(*search)
	return nil
}
