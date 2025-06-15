package cmdcopy

import (
	"passwords/data"

	"github.com/atotto/clipboard"
	"github.com/binary-soup/go-command/command"
	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type CopyCommand struct {
	command.CommandBase
}

func NewCopyCommand() CopyCommand {
	return CopyCommand{
		CommandBase: command.NewCommandBase("copy", "copy password data to the clipboard"),
	}
}

func (cmd CopyCommand) Run(args []string) error {
	name := cmd.Flags.String("name", "", "name of the vault item")
	u := cmd.Flags.Bool("u", false, "copy username")
	url := cmd.Flags.Bool("url", false, "copy url")
	p := cmd.Flags.Bool("p", true, "copy password")
	cmd.Flags.Parse(args)

	cfg, err := data.LoadConfig()
	if err != nil {
		return err
	}

	if *name == "" {
		return util.Error("name must not be empty")
	}
	filename := *name + ".json"

	password, err := cfg.Vault.LoadPassword(filename)
	if err != nil {
		return err
	}

	field := ""
	if *u {
		field = "username"
		err = cmd.copyToClipboard(password.Username)
	} else if *url {
		field = "URL"
		err = cmd.copyToClipboard(password.URL)
	} else if *p {
		field = "password"
		err = cmd.copyToClipboard(password.Password)
	}

	if err != nil {
		return util.ChainError(err, "error copying to clipboard")
	}

	style.New(style.Yellow).PrintF("%s.%s", *name, field)
	style.BoldInfo.Println(" -> Copied to Clipboard!")

	return nil
}

func (CopyCommand) copyToClipboard(text string) error {
	err := clipboard.WriteAll(text)
	if err != nil {
		return util.ChainError(err, "error copying to clipboard")
	}

	return nil
}
