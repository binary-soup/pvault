package tools

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/binary-soup/go-command/style"
)

func PromptOverwrite(title, path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return true
	}

	fmt.Printf("%s file \"%s\" exists. Overwrite [Y/n]? ", style.Bolded.Format(title), path)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}

	line := strings.TrimSpace(scanner.Text())

	if len(line) == 0 {
		return false
	} else {
		return line[0] == 'Y'
	}
}
