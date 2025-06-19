package tools

import (
	"github.com/atotto/clipboard"
	"github.com/binary-soup/go-command/util"
)

func CopyToClipboard(text string) error {
	err := clipboard.WriteAll(text)
	if err != nil {
		return util.ChainError(err, "error copying to clipboard")
	}

	return nil
}
