package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"pvault/data"
	"pvault/data/config"
	"pvault/data/vault"

	"github.com/binary-soup/go-commando/command"
	commando_config "github.com/binary-soup/go-commando/config"
	"github.com/binary-soup/go-commando/style"
	"github.com/binary-soup/go-commando/util"
)

type InitCommand struct {
	command.CommandBase
}

func NewInitCommand() InitCommand {
	return InitCommand{
		CommandBase: command.NewCommandBase("init", "initialize a new profile or upgrade an existing one"),
	}
}

func (cmd InitCommand) Run(args []string) error {
	profile := cmd.Flags.String("prof", "main", "the config profile")
	cmd.Flags.Parse(args)

	path := commando_config.GetProfilePath(DATA_DIR, *profile)

	cfg, err := cmd.initProfile(*profile, path)
	if err != nil {
		return err
	}

	err = cmd.initVault(cfg.Vault)
	if err != nil {
		return err
	}

	return nil
}

func (cmd InitCommand) initProfile(profile, path string) (*config.Config, error) {
	ok, err := data.FileExists(path)
	if err != nil {
		return nil, err
	}
	if ok {
		fmt.Printf("Profile (up to date) [%s]\n", style.Info.Format(path))
		return util.LoadJSON[config.Config]("profile config", path)
	}

	cfg := &config.Config{
		Vault: &vault.Vault{
			Path: filepath.Join(filepath.Dir(path), profile+".vault"),
		},
		Passkey: config.PasskeyConfig{
			Timeout: config.DEFAULT_PASSKEY_TIMEOUT,
		},
		Password: config.PasswordConfig{
			Lifetime: config.DEFAULT_PASSWORD_LIFETIME,
		},
	}

	os.MkdirAll(filepath.Dir(path), 0755)
	err = util.SaveJSON("profile config", cfg, path)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Created new Profile [%s]\n", style.Create.Format(path))

	return cfg, nil
}

func (cmd InitCommand) initVault(v *vault.Vault) error {
	ok, err := data.DirExists(v.Path)
	if err != nil {
		return err
	}

	if !ok {
		return cmd.createVault(v)
	}

	//TODO: check version info and upgrade if needed.

	fmt.Printf("Vault (up to date) [%s]\n", style.Info.Format(v.Path))
	return nil
}

func (cmd InitCommand) createVault(v *vault.Vault) error {
	err := v.Create()
	if err != nil {
		return err
	}

	fmt.Printf("Created new Vault [%s]\n", style.Create.Format(v.Path))
	return nil
}
