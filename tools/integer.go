package tools

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/binary-soup/go-command/alert"
	"github.com/binary-soup/go-command/style"
)

func ReadInteger(title string, min, max int) (int, error) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Enter %s: ", style.Bolded.Format(title))

		if !scanner.Scan() {
			return -1, alert.Error("error reading input")
		}

		n, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil || n < min || n > max {
			continue
		}
		return n, nil
	}
}
