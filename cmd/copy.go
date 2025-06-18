package cmd

import (
	"fmt"
	"passwords/data"
	vw "passwords/workflows/vault"

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
	search := cmd.Flags.String("s", "", "the search term")
	u := cmd.Flags.Bool("u", false, "copy username")
	url := cmd.Flags.Bool("url", false, "copy url")
	p := cmd.Flags.Bool("p", true, "copy password")
	cmd.Flags.Parse(args)

	cfg, err := data.LoadConfig()
	if err != nil {
		return err
	}

	if *search == "" {
		return util.Error("(s)earch missing or invalid")
	}

	workflow := vw.NewVaultWorkflow(cfg.Vault)

	name, err := workflow.SearchExactName(*search)
	if err != nil {
		return err
	}

	password, _, err := workflow.Decrypt(name)
	if err != nil {
		return err
	}

	field := ""
	if *u {
		field = "USERNAME"
		err = cmd.copyToClipboard(password.Username)
	} else if *url {
		field = "URL"
		err = cmd.copyToClipboard(password.URL)
	} else if *p {
		field = "PASSWORD"
		err = cmd.copyToClipboard(password.Password)
	}

	if err != nil {
		return util.ChainError(err, "error copying to clipboard")
	}

	fmt.Printf("%s.%s -> %s\n", NAME_STYLE.FormatF("\"%s\"", name), style.BoldUnderline.Format(field), style.BoldInfo.Format("Copied to Clipboard"))
	return nil
}

func (CopyCommand) copyToClipboard(text string) error {
	err := clipboard.WriteAll(text)
	if err != nil {
		return util.ChainError(err, "error copying to clipboard")
	}

	return nil
}
