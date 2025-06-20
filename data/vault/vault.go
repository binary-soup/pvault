package vault

import (
	"os"

	"github.com/binary-soup/go-command/util"
)

type Vault struct {
	Path string `json:"path"`

	index indexMap
}

func (v *Vault) Open() error {
	err := os.Mkdir(v.Path, 0755)
	if err != nil && !os.IsExist(err) {
		return util.ChainError(err, "error creating vault directory")
	}

	v.index, err = v.loadIndex()
	if err != nil {
		return util.ChainError(err, "error loading index map")
	}

	return nil
}

func (v Vault) Close() error {
	return v.saveIndex()
}
