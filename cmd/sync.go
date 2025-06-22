package cmd

import (
	syncworkflow "pvault/workflows/sync"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type SyncCommand struct {
	ConfigCommandBase
}

func NewSyncCommand() SyncCommand {
	return SyncCommand{
		ConfigCommandBase: NewConfigCommandBase("sync", "sync files between a host and client vault"),
	}
}

func (cmd SyncCommand) Run(args []string) error {
	host := cmd.Flags.Bool("host", false, "run as the host")
	port := cmd.Flags.String("port", ":9000", "port to run with/connect to")
	addr := cmd.Flags.String("addr", "", "address of the host to sync to")
	persist := cmd.Flags.Bool("persist", false, "keep the host open after syncing")
	fresh := cmd.Flags.Bool("fresh", false, "preform a fresh sync (clears the filter)")
	cmd.Flags.Parse(args)

	cfg, err := cmd.LoadConfig()
	if err != nil {
		return err
	}
	defer cfg.Vault.Close()

	if *host {
		return syncworkflow.NewHostWorkflow(cfg.Vault, *persist).Run(*port)
	}

	if *addr == "" {
		return util.Error("(addr)ess missing or empty")
	}

	if *fresh {
		cfg.Vault.Filter.Clear()
		style.Info.Println("(cleared the filter)")
	}

	return syncworkflow.NewClientWorkflow(cfg.Vault).Run(*addr, *port)
}
