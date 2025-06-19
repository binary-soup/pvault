package cmd

import (
	"fmt"
	"os"
	"passwords/data"
	"passwords/data/vault"
	vw "passwords/workflows/vault"

	"github.com/binary-soup/go-command/command"
	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type LockCommand struct {
	command.CommandBase
}

func NewLockCommand() LockCommand {
	return LockCommand{
		CommandBase: command.NewCommandBase("lock", "encrypt the provided password file"),
	}
}

func (cmd LockCommand) Run(args []string) error {
	name := cmd.Flags.String("n", "", "name of the password file")
	keep := cmd.Flags.Bool("keep", false, "keep the original password file")
	cmd.Flags.Parse(args)

	cfg, err := data.LoadConfig()
	if err != nil {
		return err
	}

	if *name == "" {
		return util.Error("(n)ame missing or invalid")
	}

	password, err := data.LoadPasswordFile(*name + ".json")
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

	cache, err := vault.LoadCacheFile(*name + ".cache.json")
	if err != nil {
		cache = &vault.Cache{}
	}

	workflow := vw.NewVaultWorkflow(cfg.Vault)
	defer cfg.Vault.Close()

	err = workflow.ChooseOrVerifyPasskey(&cache.Passkey)
	if err != nil {
		return err
	}

	err = workflow.Encrypt(password, cache.Passkey)
	if err != nil {
		return err
	}

	if !*keep {
		os.Remove(*name + ".json")
		os.Remove(*name + ".cache.json")
	}

	fmt.Printf("%s -> %s\n", NAME_STYLE.FormatF("\"%s\"", password.Name), style.BoldInfo.Format("Saved to Vault"))
	return nil
}
