package cmdworkflow

import (
	"fmt"
	cmdstyle "passwords/cmd/style"
	"passwords/data"
	"passwords/tools"
	vw "passwords/workflows/vault"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

func (cmd CmdWorkflow) RunDecrypt(cmdName string, delete bool) error {
	search := cmd.flags.String("s", "", "the search term")
	out := cmd.flags.String("o", "", "name of the out file. Defaults to search term (+.json)")
	cmd.flags.Parse(cmd.args)

	cfg, err := data.LoadConfig()
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

	if !delete {
		password.Cache = cache
	}

	err = password.SaveToFile(*out)
	if err != nil {
		return err
	}

	if delete {
		err = workflow.Delete(name)
		if err != nil {
			return err
		}
	}

	fmt.Printf("%s -> %s\n", cmdstyle.NAME_STYLE.FormatF("\"%s\"", password.Name), style.BoldInfo.FormatF("%s from Vault", cmdName))
	return nil
}
