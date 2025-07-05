package tools

import (
	"github.com/atotto/clipboard"
	"github.com/binary-soup/go-command/alert"
)

func CopyToClipboard(text string) error {
	err := clipboard.WriteAll(text)
	if err != nil {
		return alert.ChainError(err, "error copying to clipboard")
	}

	return nil
}
