package cmd

import (
	vw "pvault/workflows/vault"
)

type SearchCommand struct {
	ConfigCommandBase
}

func NewSearchCommand() SearchCommand {
	return SearchCommand{
		ConfigCommandBase: NewConfigCommandBase("search", "search the items in the vault"),
	}
}

func (cmd SearchCommand) Run(args []string) error {
	search := cmd.Flags.String("s", "", "the search term. Leave blank to list all")
	cmd.Flags.Parse(args)

	cfg, err := cmd.LoadConfig()
	if err != nil {
		return err
	}

	vw.NewVaultWorkflow(cfg.Vault).Search(*search)
	return nil
}
