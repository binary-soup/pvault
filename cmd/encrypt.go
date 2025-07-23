package cmd

import (
	"fmt"
	"os"
	"pvault/data/config"
	"pvault/data/password"
	"pvault/data/vault"
	"pvault/tools"
	vw "pvault/workflows/vault"

	"github.com/binary-soup/go-commando/alert"
	"github.com/binary-soup/go-commando/command"
	"github.com/binary-soup/go-commando/style"
)

type EncryptCommandBase struct {
	command.ConfigCommandBase[config.Config]
	new bool
}

func newEncryptCommandBase(name, desc string, new bool) EncryptCommandBase {
	return EncryptCommandBase{
		ConfigCommandBase: command.NewConfigCommandBase[config.Config](name, desc),
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

	cfg, err := cmd.LoadConfig(DATA_DIR)
	if err != nil {
		return err
	}

	if *path == "" {
		return alert.Error("(p)ath missing or invalid")
	}

	var cache *password.Cache
	if cmd.new {
		cache = &password.Cache{}
		cache.Password, err = password.LoadFile(*path)
		cache.Meta = password.NewMeta("", "")
	} else {
		cache, err = password.LoadCacheFile(*path)
	}
	if err != nil {
		return err
	}

	if !cmd.new && !cfg.Vault.Index.HasID(cache.Meta.ID) {
		return alert.ErrorF("id \"%s\" not found", cache.Meta.ID.String())
	}

	err = cache.Password.Validate()
	if err != nil {
		return alert.ChainError(err, "error validating password")
	}

	cmd.promptNewName(cfg.Vault.Index, cache.Meta)

	err = tools.PromptPasskey(&cache.Meta.Passkey)
	if err != nil {
		return err
	}

	workflow := vw.NewVaultWorkflow(cfg.Vault)
	defer cfg.Vault.Close()

	err = workflow.Encrypt(cache)
	if err != nil {
		return err
	}

	if !*keep {
		os.Remove(*path)
	}

	fmt.Printf("%s -> %s %s\n", NAME_STYLE.FormatF("\"%s\"", cache.Meta.Name), style.Info.FormatF("%s in", cmd.Name), style.BoldInfo.Format("VAULT"))
	return nil
}

func (cmd EncryptCommandBase) promptNewName(index *vault.Index, meta *password.Meta) string {
	for meta.Name == "" || index.HasName(meta.Name) {
		if meta.Name != "" {
			style.Error.Println("(name already in use)")
		}
		meta.Name = tools.PromptString(true, fmt.Sprintf("Choose New %s:", style.Bolded.Format("NAME")))
	}
	return meta.Name
}
