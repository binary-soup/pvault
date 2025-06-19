package vault

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/binary-soup/go-command/util"
)

const INDEX_FILE = "index.txt"

func (v Vault) NameExists(name string) bool {
	return v.index.Has(name)
}

func (v Vault) saveIndex() error {
	file, err := os.Create(filepath.Join(v.Path, INDEX_FILE))
	if err != nil {
		return util.ChainError(err, "error creating index file")
	}
	defer file.Close()

	for name := range v.index {
		fmt.Fprintln(file, name)
	}
	return nil
}

func (v Vault) loadIndex() (stringSet, error) {
	file, err := os.Open(filepath.Join(v.Path, INDEX_FILE))
	if os.IsNotExist(err) {
		return stringSet{}, nil
	}
	if err != nil {
		return nil, util.ChainError(err, "error opening index file")
	}
	defer file.Close()

	index := stringSet{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		index.Add(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, util.ChainError(err, "error parsing index file")
	}
	return index, nil
}
