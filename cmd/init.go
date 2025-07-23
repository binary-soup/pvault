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

	ok, err := data.PathExists(path)
	if err != nil {
		return err
	}
	if ok {
		style.Info.Println("[UP TO DATE]")
		return nil
	}

	err = cmd.initProfile(*profile, path)
	if err != nil {
		return err
	}

	fmt.Printf("Created profile %s [%s]\n", *profile, style.Create.Format(path))
	return nil
}

func (cmd InitCommand) initProfile(profile, path string) error {
	cfg := config.Config{
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
	return util.SaveJSON("profile config", &cfg, path)
}
