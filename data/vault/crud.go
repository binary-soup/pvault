package vault

import (
	"crypto/sha256"
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/binary-soup/go-command/util"
)

func (v Vault) SaveData(bytes []byte, name string) error {
	err := os.WriteFile(v.filepath(name), bytes, 0600)
	if err != nil {
		return util.ChainError(err, "error saving file to vault")
	}

	v.index.Add(name)
	return nil
}

func (v Vault) ReadData(name string) ([]byte, error) {
	bytes, err := os.ReadFile(v.filepath(name))
	if err != nil {
		return nil, util.ChainError(err, "error reading file from vault")
	}

	return bytes, nil
}

func (v Vault) DeleteData(name string) error {
	err := os.Remove(v.filepath(name))
	if err != nil {
		return util.ChainError(err, "error deleting file from vault")
	}

	v.index.Delete(name)
	return nil
}

func (v Vault) filepath(name string) string {
	hash := sha256.Sum256([]byte(name))
	return filepath.Join(v.Path, base64.RawURLEncoding.EncodeToString(hash[:])+".crypt")
}
