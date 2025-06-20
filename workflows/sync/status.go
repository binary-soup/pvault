package sw

import (
	"fmt"

	"github.com/binary-soup/go-command/style"
)

func (w SyncWorkflow) printSuccessStatus(message string) {
	fmt.Printf("  %s %s\n", style.Success.Format("[+]"), message)
}

func (w SyncWorkflow) printErrorStatus(message string) {
	fmt.Printf("  %s %s\n", style.Error.Format("[X]"), message)
}
