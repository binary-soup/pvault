package cmd

import (
	"flag"
	"pvault/data/config"
	syncworkflow "pvault/workflows/sync"

	"github.com/binary-soup/go-commando/alert"
	"github.com/binary-soup/go-commando/command"
	"github.com/binary-soup/go-commando/style"
)

type SyncCommand struct {
	command.ConfigCommandBase[config.Config]
	flags *syncFlags
}

type syncFlags struct {
	Host    *bool
	Port    *string
	Address *string

	Persist *bool
	Fresh   *bool
}

func (f *syncFlags) Set(flags *flag.FlagSet) {
	f.Host = flags.Bool("host", false, "run as the host")
	f.Port = flags.String("port", ":9000", "port to run with/connect to")
	f.Address = flags.String("addr", "", "address of the host to sync to")
	f.Persist = flags.Bool("persist", false, "keep the host open after syncing")
	f.Fresh = flags.Bool("fresh", false, "preform a fresh sync (clears the filter)")
}

func NewSyncCommand() SyncCommand {
	flags := new(syncFlags)

	return SyncCommand{
		ConfigCommandBase: command.NewConfigCommandBase[config.Config]("sync", "sync files between a host and client vault", flags),
		flags:             flags,
	}
}

func (cmd SyncCommand) Run() error {
	cfg, err := cmd.LoadConfig()
	if err != nil {
		return err
	}
	defer cfg.Vault.Close()

	if *cmd.flags.Host {
		return syncworkflow.NewHostWorkflow(cfg.Vault).Run(*cmd.flags.Port, *cmd.flags.Persist)
	}

	if *cmd.flags.Address == "" {
		return alert.Error("(addr)ess missing or empty")
	}

	if *cmd.flags.Fresh {
		cfg.Vault.Filter.Clear()
		style.Info.Println("(cleared the filter)")
	}

	return syncworkflow.NewClientWorkflow(cfg.Vault).Run(*cmd.flags.Address, *cmd.flags.Port)
}
