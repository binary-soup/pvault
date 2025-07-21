package vault

import (
	"os"

	"github.com/binary-soup/go-commando/alert"
)

type Vault struct {
	Path   string `json:"path"`
	Index  *Index
	Filter *Filter
}

func (v *Vault) Open() error {
	err := os.Mkdir(v.Path, 0755)
	if err != nil && !os.IsExist(err) {
		return alert.ChainError(err, "error creating vault directory")
	}

	v.Index, err = v.loadIndex()
	if err != nil {
		return alert.ChainError(err, "error loading index")
	}

	v.Filter, err = v.loadFilter()
	if err != nil {
		return alert.ChainError(err, "error loading filter")
	}

	return nil
}

func (v Vault) Close() {
	v.Flush()
}

func (v Vault) Flush() {
	v.saveIndex()
	v.saveFilter()
}
