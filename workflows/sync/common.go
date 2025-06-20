package syncworkflow

import (
	"fmt"
	"os"

	"github.com/binary-soup/go-command/style"
)

func hostname() string {
	name, _ := os.Hostname()
	return name
}

func printSuccessStatus(message string) {
	fmt.Printf("  %s %s\n", style.Success.Format("[+]"), message)
}

func printErrorStatus(message string) {
	fmt.Printf("  %s %s\n", style.Error.Format("[X]"), message)
}
