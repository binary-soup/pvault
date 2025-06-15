package data

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

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

func (v Vault) Delete(filename string) error {
	err := os.Remove(v.getFilepath(filename))
	if err != nil {
		return util.ChainError(err, "error deleting crypt file from vault")
	}

	return nil
}

func (v Vault) Search(substring string) ([]string, error) {
	entries, err := os.ReadDir(v.Path)
	if err != nil {
		return nil, util.ChainError(err, "error reading vault directory")
	}

	items := []string{}

	for _, entry := range entries {
		name, ok := strings.CutSuffix(entry.Name(), ".json.crypt")
		if !ok {
			continue
		}

		if substring == "" || strings.Contains(name, substring) {
			items = append(items, name)
		}
	}

	sort.Strings(items)
	return items, nil
}
