package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/binary-soup/go-command/util"
	"github.com/google/uuid"
)

const INDEX_FILE = "index.txt"

func (v Vault) NameExists(name string) bool {
	_, ok := v.index[name]
	return ok
}

func (v Vault) saveIndex(name string, id uuid.UUID) error {
	v.index[name] = id
	return v.writeIndexFile()
}

func (v Vault) deleteIndex(name string) error {
	delete(v.index, name)
	return v.writeIndexFile()
}

func (v Vault) writeIndexFile() error {
	file, err := os.Create(filepath.Join(v.Path, INDEX_FILE))
	if err != nil {
		return util.ChainError(err, "error creating index file")
	}
	defer file.Close()

	for name, id := range v.index {
		fmt.Fprintf(file, "%s:%s\n", id.String(), name)
	}
	return nil
}
