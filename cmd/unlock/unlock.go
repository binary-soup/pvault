package cmdunlock

import (
	"passwords/data"
	"path/filepath"

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
	path := cmd.Flags.String("path", "", "path to the output file. The name of the file should match the name in the vault")
	cmd.Flags.Parse(args)

	cfg, err := data.LoadConfig()
	if err != nil {
		return err
	}

	if *path == "" {
		return util.Error("path must not be empty")
	}
	filename := filepath.Base(*path)

	password, err := cfg.Vault.LoadPassword(filename)
	if err != nil {
		return err
	}

	err = password.SaveToFile(*path)
	if err != nil {
		return err
	}

	err = cfg.Vault.DeleteCrypt(filename)
	if err != nil {
		return err
	}

	style.BoldInfo.Print("Loaded from Vault -> ")
	style.New(style.Yellow).Println(*path)

	return nil
}
