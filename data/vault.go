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

func (v Vault) LoadPassword(filename string) (*Password, error) {
	bytes, err := os.ReadFile(v.getFilepath(filename))
	if err != nil {
		return nil, util.ChainError(err, "error reading crypt file from vault")
	}

	password, err := DecryptPassword(bytes)
	if err != nil {
		return nil, err
	}

	return password, nil
}

func (v Vault) SavePassword(password *Password, filename string) error {
	bytes, err := password.Encrypt()
	if err != nil {
		return err
	}

	err = os.WriteFile(v.getFilepath(filename), bytes, 0600)
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
