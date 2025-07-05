package cmd

import (
	"fmt"
	"pvault/tools"
	vw "pvault/workflows/vault"

	"github.com/binary-soup/go-command/alert"
	"github.com/binary-soup/go-command/style"
)

type CopyCommand struct {
	ConfigCommandBase
}

func NewCopyCommand() CopyCommand {
	return CopyCommand{
		ConfigCommandBase: NewConfigCommandBase("copy", "copy password data to the clipboard"),
	}
}

func (cmd CopyCommand) Run(args []string) error {
	search := cmd.Flags.String("s", "", "the search term")
	u := cmd.Flags.Bool("u", false, "copy username")
	url := cmd.Flags.Bool("url", false, "copy url")
	p := cmd.Flags.Bool("p", true, "copy password")
	cmd.Flags.Parse(args)

	cfg, err := cmd.LoadConfig()
	if err != nil {
		return err
	}

	if *search == "" {
		return alert.Error("(s)earch missing or invalid")
	}

	workflow := vw.NewVaultWorkflow(cfg.Vault)

	name, err := workflow.SearchExactName(*search)
	if err != nil {
		return err
	}

	cache, err := workflow.Decrypt(name, cfg.Passkey.Timeout)
	if err != nil {
		return err
	}

	field := ""
	if *u {
		field = "USERNAME"
		err = tools.CopyToClipboard(cache.Password.Username)
	} else if *url {
		field = "URL"
		err = tools.CopyToClipboard(cache.Password.URL)
	} else if *p {
		field = "PASSWORD"
		err = tools.CopyToClipboard(cache.Password.Password)
	}

	if err != nil {
		return err
	}

	fmt.Printf("%s.%s -> %s\n", NAME_STYLE.FormatF("\"%s\"", name), style.BoldUnderline.Format(field), style.BoldInfo.Format("Copied to Clipboard"))
	return nil
}
