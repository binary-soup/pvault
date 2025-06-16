package tools

import (
	"fmt"
	"os"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
	"golang.org/x/term"
)

func ReadPasskey(prompt string) (string, error) {
	fmt.Printf("%s %s: ", prompt, style.Bolded.Format("PASSKEY"))

	bytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", util.ChainError(err, "error reading passkey from terminal")
	}

	fmt.Println()
	return string(bytes), nil
}

func VerifyPasskey(passkey string) error {
	for {
		input, err := ReadPasskey("Verify")
		if err != nil {
			return err
		}

		if input == passkey {
			return nil
		}
	}
}

func ReadAndVerifyPasskey(prompt string) (string, error) {
	passkey, err := ReadPasskey(prompt)
	if err != nil {
		return "", err
	}

	return passkey, VerifyPasskey(passkey)
}
