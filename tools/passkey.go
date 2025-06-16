package tools

import (
	"fmt"
	"os"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
	"golang.org/x/term"
)

func ReadPasskey() (string, error) {
	fmt.Printf("Enter %s: ", style.Bolded.Format("PASSKEY"))

	bytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", util.ChainError(err, "error reading passkey from terminal")
	}

	fmt.Println()
	return string(bytes), nil
}
