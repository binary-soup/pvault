package data

import (
	"os"

	"github.com/binary-soup/go-commando/alert"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, alert.ChainError(err, "error checking file status")
	}

	return true, nil
}
