package cmd

import (
	"flag"
	"pvault/data/config"
	vw "pvault/workflows/vault"

	"github.com/binary-soup/go-commando/command"
)

type SearchCommand struct {
	command.ConfigCommandBase[config.Config]
	flags *searchFlags
}

type searchFlags struct {
	Search *string
}

func (f *searchFlags) Set(flags *flag.FlagSet) {
	f.Search = flags.String("s", "", "the search term. Leave blank to list all")
}

func NewSearchCommand() SearchCommand {
	flags := new(searchFlags)

	return SearchCommand{
		ConfigCommandBase: command.NewConfigCommandBase[config.Config]("search", "search the items in the vault", flags),
		flags:             flags,
	}
}

func (cmd SearchCommand) Run() error {
	cfg, err := cmd.LoadConfig()
	if err != nil {
		return err
	}

	vw.NewVaultWorkflow(cfg.Vault).Search(*cmd.flags.Search)
	return nil
}
