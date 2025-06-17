package cmd

import (
	"fmt"
	"passwords/data"
	"passwords/workflows"
	"strings"

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
	substring := cmd.Flags.String("s", "", "the search term. Leave blank to list all")
	cmd.Flags.Parse(args)

	cfg, err := data.LoadConfig()
	if err != nil {
		return err
	}

	items, err := workflows.SearchVault(cfg.Vault, *substring)
	if err != nil {
		return err
	}

	for _, item := range items {
		cmd.styleItem(item, *substring)
	}

	return nil
}

func (cmd SearchCommand) styleItem(item string, substring string) {
	if substring == "" {
		ITEM_STYLE.Println(item)
		return
	}

	before, after, _ := strings.Cut(item, substring)

	fmt.Printf("%s%s%s\n", ITEM_STYLE.Format(before), ITEM_HIGHLIGHT.Format(substring), ITEM_STYLE.Format(after))
}
