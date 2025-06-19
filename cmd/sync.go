package cmd

import (
	sw "passwords/workflows/sync"

	"github.com/binary-soup/go-command/command"
	"github.com/binary-soup/go-command/util"
)

type SyncCommand struct {
	command.CommandBase
}

func NewSyncCommand() SyncCommand {
	return SyncCommand{
		CommandBase: command.NewCommandBase("sync", "sync files between a host and client vault"),
	}
}

func (cmd SyncCommand) Run(args []string) error {
	host := cmd.Flags.Bool("host", false, "run as the host")
	addr := cmd.Flags.String("addr", "", "address of the host to sync to")
	cmd.Flags.Parse(args)

	workflow := sw.NewSyncWorkflow()

	if *host {
		return workflow.RunHost()
	}

	if *addr == "" {
		return util.Error("(addr)ess missing or empty")
	}

	return workflow.RunClient(*addr)
}
