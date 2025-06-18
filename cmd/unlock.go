package cmd

import (
	"fmt"
	"passwords/data"
	"passwords/tools"
	vw "passwords/workflows/vault"

	"github.com/binary-soup/go-command/command"
	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type UnlockCommand struct {
	command.CommandBase
}

func NewUnlockCommand() UnlockCommand {
	return UnlockCommand{
		CommandBase: command.NewCommandBase("unlock", "decrypt the provided crypt file"),
	}
}

func (cmd UnlockCommand) Run(args []string) error {
	search := cmd.Flags.String("s", "", "the search term")
	out := cmd.Flags.String("o", "", "name of the out file. Defaults to search term")
	cache := cmd.Flags.Bool("cache", false, "cache the passkey for faster re-lock")
	cmd.Flags.Parse(args)

	cfg, err := data.LoadConfig()
	if err != nil {
		return err
	}

	if *search == "" {
		return util.Error("(s)earch missing or invalid")
	}
	if *out == "" {
		*out = *search
	}

	workflow := vw.NewVaultWorkflow(cfg.Vault)

	name, err := workflow.SearchExactName(*search)
	if err != nil {
		return err
	}

	if ok := tools.PromptOverwrite("OUT", *out+".json"); !ok {
		return nil
	}

	password, index, err := workflow.Decrypt(name)
	if err != nil {
		return err
	}

	err = password.SaveToFile(*out + ".json")
	if err != nil {
		return err
	}

	if *cache {
		err = index.SaveToFile(*out + ".cache.json")
		if err != nil {
			return err
		}
	}

	err = workflow.Delete(name)
	if err != nil {
		return err
	}

	fmt.Printf("%s -> %s\n", NAME_STYLE.FormatF("\"%s\"", password.Name), style.BoldInfo.Format("Loaded from Vault"))
	return nil
}
