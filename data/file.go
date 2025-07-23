package data

import (
	"os"

	"github.com/binary-soup/go-commando/alert"
)

func pathExists(path string) (os.FileInfo, bool, error) {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return stat, false, nil
	}
	if err != nil {
		return stat, false, alert.ChainError(err, "error checking file status")
	}

	return stat, true, nil
}

func FileExists(path string) (bool, error) {
	stat, ok, err := pathExists(path)
	if !ok || err != nil {
		return false, err
	}

	if stat.IsDir() {
		return true, alert.ErrorF("file \"%s\" is a directory", path)
	}
	return true, nil
}

func DirExists(path string) (bool, error) {
	stat, ok, err := pathExists(path)
	if !ok || err != nil {
		return false, err
	}

	if !stat.IsDir() {
		return true, alert.ErrorF("dir \"%s\" is a not directory", path)
	}
	return true, nil
}
