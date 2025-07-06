package cmd

import (
	"fmt"
	"math"
	"pvault/data/config"
	"pvault/tools"
	vw "pvault/workflows/vault"

	"github.com/binary-soup/go-command/alert"
	"github.com/binary-soup/go-command/command"
	"github.com/binary-soup/go-command/style"
)

type CopyCommand struct {
	command.ConfigCommandBase[config.Config]
}

func NewCopyCommand() CopyCommand {
	return CopyCommand{
		ConfigCommandBase: command.NewConfigCommandBase[config.Config]("copy", "copy password data to the clipboard"),
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

	NAME_STYLE.Print(name)

	if *u {
		err = cmd.copyToClipboard(cfg, "USERNAME", cache.Password.Username)
	} else if *url {
		err = cmd.copyToClipboard(cfg, "URL", cache.Password.URL)
	} else if *p {
		err = cmd.copyToClipboard(cfg, "PASSWORD", cache.Password.Password)
	}

	if err != nil {
		return err
	}

	return nil
}

func (cmd CopyCommand) copyToClipboard(cfg *config.Config, field, text string) error {
	fmt.Printf(".%s -> %s\n", style.Bolded.Format(field), style.BoldInfo.Format("Copied to Clipboard"))

	ch, err := tools.TempCopyToClipboard(text, cfg.Password.Lifetime, "REDACTED")
	if ch == nil || err != nil {
		return err
	}

	for range int(math.Floor(float64(cfg.Password.Lifetime))) {
		tools.Timeout(1)
		style.Bolded.Print("-")
	}

	<-ch
	style.Delete.Println("\nREDACTED")

	return nil
}
