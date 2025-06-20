package cmd

import (
	"pvault/data"
	vw "pvault/workflows/vault"

	"github.com/binary-soup/go-command/command"
)

type SearchCommand struct {
	command.CommandBase
}

func NewSearchCommand() SearchCommand {
	return SearchCommand{
		CommandBase: command.NewCommandBase("search", "search the items in the vault"),
	}
}

func (cmd SearchCommand) Run(args []string) error {
	search := cmd.Flags.String("s", "", "the search term. Leave blank to list all")
	cmd.Flags.Parse(args)

	cfg, err := data.LoadConfig()
	if err != nil {
		return err
	}

	vw.NewVaultWorkflow(cfg.Vault).Search(*search)
	return nil
}
