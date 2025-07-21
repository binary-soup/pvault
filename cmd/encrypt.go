package cmd

import (
	"flag"
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

type encryptCommandBase struct {
	command.ConfigCommandBase[config.Config]
	flags  *encryptFlags
	create bool
}

type encryptFlags struct {
	Path *string
	Keep *bool
}

func (f *encryptFlags) Set(flags *flag.FlagSet) {
	f.Path = flags.String("p", "", "path to the password file")
	f.Keep = flags.Bool("keep", false, "keep the original password file")
}

func newEncryptCommandBase(name, desc string, create bool) encryptCommandBase {
	flags := new(encryptFlags)

	return encryptCommandBase{
		ConfigCommandBase: command.NewConfigCommandBase[config.Config](name, desc, flags),
		flags:             flags,
		create:            create,
	}
}

//#######################

type RelockCommand struct {
	encryptCommandBase
}

func NewRelockCommand() RelockCommand {
	return RelockCommand{
		encryptCommandBase: newEncryptCommandBase("relock", "re-encrypt a file back into the vault", false),
	}
}

//#######################

type StashCommand struct {
	encryptCommandBase
}

func NewStashCommand() StashCommand {
	return StashCommand{
		encryptCommandBase: newEncryptCommandBase("stash", "encrypt and stash a new file in the vault", true),
	}
}

//#######################

func (cmd encryptCommandBase) Run() error {
	cfg, err := cmd.LoadConfig()
	if err != nil {
		return err
	}

	if *cmd.flags.Path == "" {
		return alert.Error("(p)ath missing or invalid")
	}

	var cache password.Cache
	if cmd.create {
		cache.Password, err = password.LoadFile(*cmd.flags.Path)
		cache.Meta = password.NewMeta("", "")
	} else {
		cache, err = password.LoadCacheFile(*cmd.flags.Path)
	}
	if err != nil {
		return err
	}

	if !cmd.create && !cfg.Vault.Index.HasID(cache.Meta.ID) {
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

	if !*cmd.flags.Keep {
		os.Remove(*cmd.flags.Path)
	}

	fmt.Printf("%s -> %s %s\n", NAME_STYLE.FormatF("\"%s\"", cache.Meta.Name), style.Info.FormatF("%s in", cmd.Name), style.BoldInfo.Format("VAULT"))
	return nil
}

func (cmd encryptCommandBase) promptNewName(index *vault.Index, meta password.Meta) string {
	for meta.Name == "" || index.HasName(meta.Name) {
		if meta.Name != "" {
			style.Error.Println("(name already in use)")
		}
		meta.Name = tools.PromptString(true, fmt.Sprintf("Choose New %s:", style.Bolded.Format("NAME")))
	}
	return meta.Name
}
