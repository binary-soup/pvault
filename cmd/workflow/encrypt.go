package cmdworkflow

import (
	"fmt"
	"os"
	cmdstyle "passwords/cmd/style"
	"passwords/data"
	vw "passwords/workflows/vault"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

func (cmd CmdWorkflow) RunEncrypt(cmdName string, new bool) error {
	path := cmd.flags.String("p", "", "path to the password file")
	keep := cmd.flags.Bool("keep", false, "keep the original password file")
	cmd.flags.Parse(cmd.args)

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

	var cache *data.PasswordCache

	if new {
		if cfg.Vault.Index.NameExists(password.Name) {
			return util.Error(fmt.Sprintf("name \"%s\" already exists", password.Name))
		}
		cache = data.NewPasswordCache("")
	} else {
		cache = password.Cache
		if cache == nil {
			return util.Error("cache missing from password file")
		}
		password.Cache = nil

		if !cfg.Vault.Index.IdExists(cache.ID) {
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

	fmt.Printf("%s -> %s\n", cmdstyle.NAME_STYLE.FormatF("\"%s\"", password.Name), style.BoldInfo.FormatF("%s in Vault", cmdName))
	return nil
}
