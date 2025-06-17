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

func (v *Vault) Init() error {
	err := os.Mkdir(v.Path, 0755)
	if err != nil && !os.IsExist(err) {
		return util.ChainError(err, "error creating vault directory")
	}

	return nil
}

func (v Vault) NewIndex() *Index {
	return &Index{}
}

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

func (v Vault) Search(substring, ext string) ([]string, error) {
	entries, err := os.ReadDir(v.Path)
	if err != nil {
		return nil, util.ChainError(err, "error reading vault directory")
	}

	items := []string{}

	for _, entry := range entries {
		name, ok := strings.CutSuffix(entry.Name(), ext)
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
