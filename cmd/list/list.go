package cmdlist

import (
	"passwords/data"

	"github.com/binary-soup/go-command/command"
	"github.com/binary-soup/go-command/style"
)

type ListCommand struct {
	command.CommandBase
}

func NewListCommand() ListCommand {
	return ListCommand{
		CommandBase: command.NewCommandBase("list", "list the items in the vault"),
	}
}

func (cmd ListCommand) Run(args []string) error {
	substring := cmd.Flags.String("s", "", "the search term. Leave blank for no search")
	cmd.Flags.Parse(args)

	cfg, err := data.LoadConfig()
	if err != nil {
		return err
	}

	items, err := cfg.Vault.Search(*substring)
	if err != nil {
		return err
	}

	style.BoldUnderline.PrintF("%d Items\n", len(items))
	for _, item := range items {
		style.New(style.Yellow).Println(item)
	}

	return nil
}
