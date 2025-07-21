package cmd

import (
	"flag"
	"fmt"
	"math"
	"pvault/data/config"
	"pvault/tools"
	vw "pvault/workflows/vault"

	"github.com/binary-soup/go-commando/alert"
	"github.com/binary-soup/go-commando/command"
	"github.com/binary-soup/go-commando/style"
)

type CopyCommand struct {
	command.ConfigCommandBase[config.Config]
	flags *copyFlags
}

type copyFlags struct {
	Search *string

	CopyUsername *bool
	CopyURL      *bool
}

func (f *copyFlags) Set(flags *flag.FlagSet) {
	f.Search = flags.String("s", "", "the search term")
	f.CopyUsername = flags.Bool("u", false, "copy username")
	f.CopyURL = flags.Bool("url", false, "copy url")
}

func NewCopyCommand() CopyCommand {
	flags := new(copyFlags)

	return CopyCommand{
		ConfigCommandBase: command.NewConfigCommandBase[config.Config]("copy", "copy password data to the clipboard", flags),
		flags:             flags,
	}
}

func (cmd CopyCommand) Run() error {
	cfg, err := cmd.LoadConfig()
	if err != nil {
		return err
	}

	if *cmd.flags.Search == "" {
		return alert.Error("(s)earch missing or invalid")
	}

	workflow := vw.NewVaultWorkflow(cfg.Vault)

	name, err := workflow.SearchExactName(*cmd.flags.Search)
	if err != nil {
		return err
	}

	cache, err := workflow.Decrypt(name, cfg.Passkey.Timeout)
	if err != nil {
		return err
	}

	NAME_STYLE.Print(name)

	if *cmd.flags.CopyUsername {
		err = cmd.copyToClipboard(cfg, "USERNAME", cache.Password.Username)
	} else if *cmd.flags.CopyURL {
		err = cmd.copyToClipboard(cfg, "URL", cache.Password.URL)
	} else { // copy password
		err = cmd.copyToClipboard(cfg, "PASSWORD", cache.Password.Password)
	}

	if err != nil {
		return err
	}

	return nil
}

func (cmd CopyCommand) copyToClipboard(cfg config.Config, field, text string) error {
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
