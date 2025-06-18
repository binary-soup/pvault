package vault

import (
	"os"

	"github.com/binary-soup/go-command/util"
)

type Vault struct {
	Path string `json:"path"`

	indexMap map[string]uint
}

func (v *Vault) Init() error {
	err := os.Mkdir(v.Path, 0755)
	if err != nil && !os.IsExist(err) {
		return util.ChainError(err, "error creating vault directory")
	}

	v.indexMap = map[string]uint{
		"Test One": 1,
	}

	return nil
}

func (v Vault) NewIndex() *Index {
	return &Index{}
}
