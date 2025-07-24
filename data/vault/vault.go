package vault

import (
	"os"

	"github.com/binary-soup/go-commando/alert"
)

type Vault struct {
	Path   string  `json:"path"`
	Index  *Index  `json:"-"`
	Filter *Filter `json:"-"`
}

func (v Vault) Create() error {
	err := os.Mkdir(v.Path, 0755)
	if err != nil {
		return alert.ChainError(err, "error creating vault directory")
	}

	err = v.writeVersion()
	if err != nil {
		return err
	}

	return nil
}

func (v *Vault) Open() error {
	var err error

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
