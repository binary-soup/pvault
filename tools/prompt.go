package tools

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/binary-soup/go-commando/style"
)

func PromptAccept(prompt string, options []byte) int {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s ", prompt)

		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 {
			continue
		}

		for i, char := range options {
			if line[0] == char {
				return i
			}
		}
	}
}

func PromptOverwrite(title, path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return true
	}

	res := PromptAccept(fmt.Sprintf("%s file \"%s\" exists. Overwrite [Y/n]?", style.Bolded.Format(title), path), []byte("Yn"))
	return res == 0
}

func PromptString(allowEmpty bool, prompt string) string {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s ", prompt)

		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())

		if allowEmpty || line != "" {
			return line
		}
	}
}

func PromptPasskey(passkey *string) error {
	var err error

	if *passkey == "" {
		*passkey, err = ReadAndVerifyPasskey("Choose New")
	} else {
		err = VerifyPasskey(*passkey)
	}

	return err
}
