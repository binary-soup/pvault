package sync

import (
	"fmt"
	"log"
	"os"

	"github.com/binary-soup/go-commando/style"
)

const LOG_INDENT = "  "

type Logger struct {
	log *log.Logger
}

func NewSuccessLogger() Logger {
	return newLogger(style.Success.Format("[+]"))
}

func NewErrorLogger() Logger {
	return newLogger(style.Error.Format("[X]"))
}

func newLogger(icon string) Logger {
	return Logger{
		log: log.New(os.Stdout, fmt.Sprintf("  %s ", icon), log.Ltime),
	}
}

func (log Logger) Log(v ...any) {
	log.log.Println(v...)
}

func (log Logger) LogF(format string, v ...any) {
	log.log.Printf(format+"\n", v...)
}
