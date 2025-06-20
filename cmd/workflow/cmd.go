package cmdworkflow

import "flag"

type CmdWorkflow struct {
	flags *flag.FlagSet
	args  []string
}

func NewCmdWorkflow(flags *flag.FlagSet, args []string) CmdWorkflow {
	return CmdWorkflow{
		flags: flags,
		args:  args,
	}
}
