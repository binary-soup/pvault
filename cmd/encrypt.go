package cmd

import (
	"fmt"
	"os"
	"pvault/data"
	vw "pvault/workflows/vault"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type EncryptCommandBase struct {
	ConfigCommandBase
	new bool
}

func newEncryptCommandBase(name, desc string, new bool) EncryptCommandBase {
	return EncryptCommandBase{
		ConfigCommandBase: NewConfigCommandBase(name, desc),
		new:               new,
	}
}

//#######################

type RelockCommand struct {
	EncryptCommandBase
}

func NewRelockCommand() RelockCommand {
	return RelockCommand{
		EncryptCommandBase: newEncryptCommandBase("relock", "re-encrypt a file back into the vault", false),
	}
}

//#######################

type StashCommand struct {
	EncryptCommandBase
}

func NewStashCommand() StashCommand {
	return StashCommand{
		EncryptCommandBase: newEncryptCommandBase("stash", "encrypt and stash a new file in the vault", true),
	}
}

//#######################

func (cmd EncryptCommandBase) Run(args []string) error {
	path := cmd.Flags.String("p", "", "path to the password file")
	keep := cmd.Flags.Bool("keep", false, "keep the original password file")
	cmd.Flags.Parse(args)

	cfg, err := cmd.LoadConfig()
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

	var cache *data.PasswordCache

	if cmd.new {
		if cfg.Vault.Index.HasName(password.Name) {
			return util.Error(fmt.Sprintf("name \"%s\" already exists", password.Name))
		}
		cache = data.NewPasswordCache("")
	} else {
		cache = password.Cache
		if cache == nil {
			return util.Error("cache missing from password file")
		}
		password.Cache = nil

		if !cfg.Vault.Index.HasID(cache.ID) {
			return util.Error(fmt.Sprintf("id \"%s\" not found", cache.ID.String()))
		}
	}

	workflow := vw.NewVaultWorkflow(cfg.Vault)
	defer cfg.Vault.Close()

	err = workflow.ChooseOrVerifyPasskey(&cache.Passkey)
	if err != nil {
		return err
	}

	err = workflow.Encrypt(password, cache)
	if err != nil {
		return err
	}

	if !*keep {
		os.Remove(*path)
	}

	fmt.Printf("%s -> %s\n", NAME_STYLE.FormatF("\"%s\"", password.Name), style.BoldInfo.FormatF("%s in Vault", cmd.Name))
	return nil
}
