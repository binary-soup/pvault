package cmd

import (
	"fmt"
	"pvault/tools"
	vw "pvault/workflows/vault"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type DecryptCommandBase struct {
	ConfigCommandBase
	remove bool
}

func newDecryptCommandBase(name, desc string, remove bool) DecryptCommandBase {
	return DecryptCommandBase{
		ConfigCommandBase: NewConfigCommandBase(name, desc),
		remove:            remove,
	}
}

//#######################

type UnlockCommand struct {
	DecryptCommandBase
}

func NewUnlockCommand() UnlockCommand {
	return UnlockCommand{
		DecryptCommandBase: newDecryptCommandBase("unlock", "temporarily decrypt a file from the vault", false),
	}
}

//#######################

type WithdrawCommand struct {
	DecryptCommandBase
}

func NewWithdrawCommand() WithdrawCommand {
	return WithdrawCommand{
		DecryptCommandBase: newDecryptCommandBase("withdraw", "decrypt and remove a file from the vault", true),
	}
}

//#######################

func (cmd DecryptCommandBase) Run(args []string) error {
	search := cmd.Flags.String("s", "", "the search term")
	out := cmd.Flags.String("o", "", "name of the out file. Defaults to search term (+.json)")
	cmd.Flags.Parse(args)

	cfg, err := cmd.LoadConfig()
	if err != nil {
		return err
	}

	if *search == "" {
		return util.Error("(s)earch missing or invalid")
	}
	if *out == "" {
		*out = *search + ".json"
	}

	workflow := vw.NewVaultWorkflow(cfg.Vault)
	defer cfg.Vault.Close()

	name, err := workflow.SearchExactName(*search)
	if err != nil {
		return err
	}

	if ok := tools.PromptOverwrite("OUT", *out); !ok {
		return nil
	}

	password, cache, err := workflow.Decrypt(name, cfg.Passkey.Timeout)
	if err != nil {
		return err
	}

	if !cmd.remove {
		password.Cache = cache
	}

	err = password.SaveToFile(*out)
	if err != nil {
		return err
	}

	if cmd.remove {
		err = workflow.Delete(name)
		if err != nil {
			return err
		}
	}

	fmt.Printf("%s -> %s\n", NAME_STYLE.FormatF("\"%s\"", password.Name), style.BoldInfo.FormatF("%s from Vault", cmd.Name))
	return nil
}
