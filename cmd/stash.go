package cmd

import (
	"fmt"
	"os"
	"passwords/data"
	vw "passwords/workflows/vault"

	"github.com/binary-soup/go-command/command"
	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type StashCommand struct {
	command.CommandBase
}

func NewStashCommand() StashCommand {
	return StashCommand{
		CommandBase: command.NewCommandBase("stash", "encrypt and stash in the vault"),
	}
}

func (cmd StashCommand) Run(args []string) error {
	path := cmd.Flags.String("p", "", "path to the password file")
	keep := cmd.Flags.Bool("keep", false, "keep the original password file")
	cmd.Flags.Parse(args)

	cfg, err := data.LoadConfig()
	if err != nil {
		return err
	}

	if *path == "" {
		return util.Error("(p)ath missing or invalid")
	}

	password, err := data.LoadPasswordFile(*path)
	if err != nil {
		return err
	}

	err = password.Validate()
	if err != nil {
		return util.ChainError(err, "error validating password")
	}

	if cfg.Vault.NameExists(password.Name) {
		return util.Error(fmt.Sprintf("name \"%s\" already exists", password.Name))
	}

	workflow := vw.NewVaultWorkflow(cfg.Vault)
	defer cfg.Vault.Close()

	var passkey string
	err = workflow.ChooseOrVerifyPasskey(&passkey)
	if err != nil {
		return err
	}

	err = workflow.Encrypt(password, passkey)
	if err != nil {
		return err
	}

	if !*keep {
		os.Remove(*path)
	}

	fmt.Printf("%s -> %s\n", NAME_STYLE.FormatF("\"%s\"", password.Name), style.BoldInfo.Format("Stashed in Vault"))
	return nil
}
