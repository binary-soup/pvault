package cmd

import "github.com/binary-soup/go-command/command"

type ImportCommand struct {
	command.CommandBase
}

func NewImportCommand() ImportCommand {
	return ImportCommand{
		CommandBase: command.NewCommandBase("import", "import many passwords from CSV [name|password|username|url]. All items will use the same passkey"),
	}
}

func (cmd ImportCommand) Run(args []string) error {
	return nil
}
