package cmd

import (
	"pvault/data"

	"github.com/binary-soup/go-command/command"
)

type ConfigCommandBase struct {
	command.CommandBase
	config *string
}

func NewConfigCommandBase(name, desc string) ConfigCommandBase {
	cmd := ConfigCommandBase{
		CommandBase: command.NewCommandBase(name, desc),
	}

	cmd.config = cmd.Flags.String("config", "", "path to a custom config file")
	return cmd
}

func (cmd ConfigCommandBase) LoadConfig() (*data.Config, error) {
	if *cmd.config != "" {
		return data.LoadCustomConfig(*cmd.config)
	}
	return data.LoadBaseConfig()
}
