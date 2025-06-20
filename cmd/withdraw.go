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

type WithdrawCommand struct {
	command.CommandBase
}

func NewWithdrawCommand() WithdrawCommand {
	return WithdrawCommand{
		CommandBase: command.NewCommandBase("withdraw", "remove and decrypt from the vault"),
	}
}

func (cmd WithdrawCommand) Run(args []string) error {
	search := cmd.Flags.String("s", "", "the search term")
	out := cmd.Flags.String("o", "", "name of the out file. Defaults to search term (+.json)")
	cmd.Flags.Parse(args)

	cfg, err := data.LoadConfig()
	if err != nil {
		return err
	}

	if *search == "" {
		return util.Error("(s)earch missing or invalid")
	}
	if *out == "" {
		*out = *search + ".json"
	}

	workflow := vw.NewVaultWorkflow(cfg.Vault)
	defer cfg.Vault.Close()

	name, err := workflow.SearchExactName(*search)
	if err != nil {
		return err
	}

	if ok := tools.PromptOverwrite("OUT", *out); !ok {
		return nil
	}

	password, _, err := workflow.Decrypt(name, cfg.Passkey.Timeout)
	if err != nil {
		return err
	}

	err = password.SaveToFile(*out)
	if err != nil {
		return err
	}

	err = workflow.Delete(name)
	if err != nil {
		return err
	}

	fmt.Printf("%s -> %s\n", NAME_STYLE.FormatF("\"%s\"", password.Name), style.BoldInfo.Format("Withdrawn from Vault"))
	return nil
}
