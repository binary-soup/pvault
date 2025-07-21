package tools

import (
	"github.com/atotto/clipboard"
	"github.com/binary-soup/go-commando/alert"
)

func CopyToClipboard(text string) error {
	_, err := TempCopyToClipboard(text, 0, "")
	return err
}

func TempCopyToClipboard(text string, lifetime float32, redactedText string) (chan struct{}, error) {
	err := clipboard.WriteAll(text)
	if err != nil {
		return nil, alert.ChainError(err, "error copying to clipboard")
	}

	if lifetime <= 0 {
		return nil, nil
	}
	ch := make(chan struct{})

	go func() {
		Timeout(lifetime)
		clipboard.WriteAll(redactedText)
		ch <- struct{}{}
	}()

	return ch, nil
}
