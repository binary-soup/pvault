package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"pvault/data"
	"pvault/data/config"
	"pvault/data/vault"
	"pvault/data/version"

	"github.com/binary-soup/go-commando/alert"
	"github.com/binary-soup/go-commando/command"
	"github.com/binary-soup/go-commando/style"
	"github.com/binary-soup/go-commando/util"
)

type InitCommand struct {
	command.ConfigCommandBase[config.Config]
}

func NewInitCommand() InitCommand {
	return InitCommand{
		ConfigCommandBase: command.NewConfigCommandBase[config.Config]("init", "initialize a new vault or upgrade an existing one"),
	}
}

func (cmd InitCommand) Run(args []string) error {
	cmd.Flags.Parse(args)

	if !cmd.UsingConfig() {
		err := cmd.initProfile()
		if err != nil {
			return err
		}
	}

	cfg, err := cmd.LoadConfig(DATA_DIR)
	if err != nil {
		return err
	}

	err = cmd.initVault(cfg.Vault)
	if err != nil {
		return err
	}

	return nil
}

func (cmd InitCommand) initProfile() error {
	path := cmd.GetConfigPath(DATA_DIR)

	ok, err := data.FileExists(path)
	if err != nil {
		return err
	}
	if ok {
		fmt.Printf("Using existing Profile %s [%s]\n", style.Bolded.Format(*cmd.Profile), style.Info.Format(path))
		return nil
	}

	cfg := config.Config{
		Vault: &vault.Vault{
			Path: filepath.Join(filepath.Dir(path), *cmd.Profile+".vault"),
		},
		Passkey: config.PasskeyConfig{
			Timeout: config.DEFAULT_PASSKEY_TIMEOUT,
		},
		Password: config.PasswordConfig{
			Lifetime: config.DEFAULT_PASSWORD_LIFETIME,
		},
	}

	os.MkdirAll(filepath.Dir(path), 0755)
	err = util.SaveJSON("profile config", &cfg, path)
	if err != nil {
		return err
	}
	fmt.Printf("Created new Profile %s [%s]\n", style.Bolded.Format(*cmd.Profile), style.Create.Format(path))

	return nil
}

func (cmd InitCommand) initVault(v *vault.Vault) error {
	ok, err := data.DirExists(v.Path)
	if err != nil {
		return err
	}

	if !ok {
		return cmd.createVault(v)
	} else {
		return cmd.upgradeVault(v)
	}
}

func (cmd InitCommand) createVault(v *vault.Vault) error {
	err := v.Create()
	if err != nil {
		return err
	}

	fmt.Printf("Created new Vault@%s [%s]\n", style.Bolded.FormatF("v%d", version.VAULT), style.Create.Format(v.Path))
	return nil
}

func (cmd InitCommand) upgradeVault(v *vault.Vault) error {
	ver, err := v.ReadVersion()
	if err != nil {
		return err
	}

	if ver == 0 || ver > version.VAULT {
		return alert.ErrorF("unsupported vault version %d", ver)
	}

	fmt.Printf("Vault@%s (up to date) [%s]\n", style.Bolded.FormatF("v%d", ver), style.Info.Format(v.Path))
	return nil
}
