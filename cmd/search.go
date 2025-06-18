package cmd

import (
	"fmt"
	"passwords/data"
	"passwords/data/vault"

	"github.com/binary-soup/go-command/command"
	"github.com/binary-soup/go-command/style"
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

	for _, item := range cfg.Vault.Search(*substring) {
		cmd.styleItem(item)
	}

	return nil
}

func (cmd SearchCommand) styleItem(item vault.SearchItem) {
	fmt.Printf("%s %s%s%s\n",
		style.BoldInfo.FormatF("[%d]", item.ID),
		ITEM_STYLE.Format(item.Name[:item.MatchStart]),
		ITEM_HIGHLIGHT.Format(item.Name[item.MatchStart:item.MatchEnd]),
		ITEM_STYLE.Format(item.Name[item.MatchEnd:]),
	)
}
