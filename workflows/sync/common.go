package syncworkflow

import (
	"os"
	"pvault/tools/sync"
)

var Success = sync.NewSuccessLogger()
var Error = sync.NewErrorLogger()

func Hostname() string {
	name, _ := os.Hostname()
	return name
}
