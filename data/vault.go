package data

import (
	"os"
	"path/filepath"

	"github.com/binary-soup/go-command/util"
)

type Vault struct {
	Path string `json:"path"`
}

func (v Vault) getFilepath(filename string) string {
	return filepath.Join(v.Path, filename) + ".crypt"
}

func (v Vault) LoadCrypt(filename string) ([]byte, error) {
	bytes, err := os.ReadFile(v.getFilepath(filename))
	if err != nil {
		return nil, util.ChainError(err, "error reading crypt file from vault")
	}

	return bytes, nil
}

func (v Vault) SaveCrypt(bytes []byte, filename string) error {
	err := os.WriteFile(v.getFilepath(filename), bytes, 0600)
	if err != nil {
		return util.ChainError(err, "error saving crypt file to vault")
	}

	return nil
}

func (v Vault) DeleteCrypt(filename string) error {
	err := os.Remove(v.getFilepath(filename))
	if err != nil {
		return util.ChainError(err, "error deleting crypt file from vault")
	}

	return nil
}
