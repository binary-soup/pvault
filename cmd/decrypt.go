package cmd

import (
	"flag"
	"fmt"
	"pvault/data/config"
	"pvault/tools"
	vw "pvault/workflows/vault"

	"github.com/binary-soup/go-commando/alert"
	"github.com/binary-soup/go-commando/command"
	"github.com/binary-soup/go-commando/style"
)

type decryptCommandBase struct {
	command.ConfigCommandBase[config.Config]
	flags  *decryptFlags
	remove bool
	delete bool
}

type decryptFlags struct {
	Search  *string
	OutPath *string

	delete bool
}

func (f *decryptFlags) Set(flags *flag.FlagSet) {
	f.Search = flags.String("s", "", "the search term")
	if !f.delete {
		f.OutPath = flags.String("o", "", "name of the out file. Defaults to search term (+.json)")
	}
}

func newDecryptCommandBase(name, desc string, remove, delete bool) decryptCommandBase {
	flags := &decryptFlags{
		delete: delete,
	}

	return decryptCommandBase{
		ConfigCommandBase: command.NewConfigCommandBase[config.Config](name, desc, flags),
		flags:             flags,
		remove:            remove,
		delete:            delete,
	}
}

//#######################

type UnlockCommand struct {
	decryptCommandBase
}

func NewUnlockCommand() UnlockCommand {
	return UnlockCommand{
		decryptCommandBase: newDecryptCommandBase("unlock", "temporarily decrypt a file from the vault", false, false),
	}
}

//#######################

type WithdrawCommand struct {
	decryptCommandBase
}

func NewWithdrawCommand() WithdrawCommand {
	return WithdrawCommand{
		decryptCommandBase: newDecryptCommandBase("withdraw", "decrypt and remove a file from the vault", true, false),
	}
}

//#######################

type DeleteCommand struct {
	decryptCommandBase
}

func NewDeleteCommand() DeleteCommand {
	return DeleteCommand{
		decryptCommandBase: newDecryptCommandBase("delete", "delete a file from the vault", true, true),
	}
}

//#######################

func (cmd decryptCommandBase) Run() error {
	cfg, err := cmd.LoadConfig()
	if err != nil {
		return err
	}

	if *cmd.flags.Search == "" {
		return alert.Error("(s)earch missing or invalid")
	}
	if *cmd.flags.OutPath == "" {
		*cmd.flags.OutPath = *cmd.flags.Search + ".json"
	}

	workflow := vw.NewVaultWorkflow(cfg.Vault)
	defer cfg.Vault.Close()

	name, err := workflow.SearchExactName(*cmd.flags.Search)
	if err != nil {
		return err
	}

	if !cmd.delete {
		if ok := tools.PromptOverwrite("OUT", *cmd.flags.OutPath); !ok {
			return nil
		}
	}

	cache, err := workflow.Decrypt(name, cfg.Passkey.Timeout)
	if err != nil {
		return err
	}

	if !cmd.delete {
		if cmd.remove {
			err = cache.Password.SaveToFile(*cmd.flags.OutPath)
		} else {
			err = cache.SaveToFile(*cmd.flags.OutPath)
		}
	}
	if err != nil {
		return err
	}

	if cmd.remove {
		err = workflow.Delete(name)
		if err != nil {
			return err
		}
	}

	fmt.Printf("%s -> %s %s\n", NAME_STYLE.FormatF("\"%s\"", name), style.Info.FormatF("%s from", cmd.Name), style.BoldInfo.Format("VAULT"))
	return nil
}
