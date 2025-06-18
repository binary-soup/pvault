package vault

import (
	"os"
	"path/filepath"

	"github.com/binary-soup/go-command/util"
)

func (v Vault) getFilepath(filename string) string {
	return filepath.Join(v.Path, filename)
}

func (v Vault) SaveData(bytes []byte, filename string) error {
	err := os.WriteFile(v.getFilepath(filename), bytes, 0600)
	if err != nil {
		return util.ChainError(err, "error saving file to vault")
	}

	return nil
}

func (v Vault) LoadData(filename string) ([]byte, error) {
	bytes, err := os.ReadFile(v.getFilepath(filename))
	if err != nil {
		return nil, util.ChainError(err, "error reading file from vault")
	}

	return bytes, err
}

func (v Vault) Delete(filename string) error {
	err := os.Remove(v.getFilepath(filename))
	if err != nil {
		return util.ChainError(err, "error deleting file from vault")
	}

	return nil
}
