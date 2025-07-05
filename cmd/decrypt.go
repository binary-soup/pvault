package cmd

import (
	"fmt"
	"pvault/data/config"
	"pvault/tools"
	vw "pvault/workflows/vault"

	"github.com/binary-soup/go-command/alert"
	"github.com/binary-soup/go-command/command"
	"github.com/binary-soup/go-command/style"
)

type DecryptCommandBase struct {
	command.ConfigCommandBase[config.Config]
	remove bool
	delete bool
}

func newDecryptCommandBase(name, desc string, remove, delete bool) DecryptCommandBase {
	return DecryptCommandBase{
		ConfigCommandBase: command.NewConfigCommandBase[config.Config](name, desc),
		remove:            remove,
		delete:            delete,
	}
}

//#######################

type UnlockCommand struct {
	DecryptCommandBase
}

func NewUnlockCommand() UnlockCommand {
	return UnlockCommand{
		DecryptCommandBase: newDecryptCommandBase("unlock", "temporarily decrypt a file from the vault", false, false),
	}
}

//#######################

type WithdrawCommand struct {
	DecryptCommandBase
}

func NewWithdrawCommand() WithdrawCommand {
	return WithdrawCommand{
		DecryptCommandBase: newDecryptCommandBase("withdraw", "decrypt and remove a file from the vault", true, false),
	}
}

//#######################

type DeleteCommand struct {
	DecryptCommandBase
}

func NewDeleteCommand() DeleteCommand {
	return DeleteCommand{
		DecryptCommandBase: newDecryptCommandBase("delete", "delete a file from the vault", true, true),
	}
}

//#######################

func (cmd DecryptCommandBase) Run(args []string) error {
	var out *string

	search := cmd.Flags.String("s", "", "the search term")
	if !cmd.delete {
		out = cmd.Flags.String("o", "", "name of the out file. Defaults to search term (+.json)")
	}
	cmd.Flags.Parse(args)

	cfg, err := cmd.LoadConfig()
	if err != nil {
		return err
	}

	if *search == "" {
		return alert.Error("(s)earch missing or invalid")
	}
	if !cmd.delete && *out == "" {
		*out = *search + ".json"
	}

	workflow := vw.NewVaultWorkflow(cfg.Vault)
	defer cfg.Vault.Close()

	name, err := workflow.SearchExactName(*search)
	if err != nil {
		return err
	}

	if !cmd.delete {
		if ok := tools.PromptOverwrite("OUT", *out); !ok {
			return nil
		}
	}

	cache, err := workflow.Decrypt(name, cfg.Passkey.Timeout)
	if err != nil {
		return err
	}

	if !cmd.delete {
		if cmd.remove {
			err = cache.Password.SaveToFile(*out)
		} else {
			err = cache.SaveToFile(*out)
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
