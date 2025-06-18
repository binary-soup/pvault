package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/binary-soup/go-command/util"
	"github.com/google/uuid"
)

func (v Vault) SaveData(bytes []byte, id uuid.UUID, name string) error {
	err := os.WriteFile(v.filepath(id), bytes, 0600)
	if err != nil {
		return util.ChainError(err, "error saving file to vault")
	}

	err = v.saveIndex(name, id)
	if err != nil {
		return util.ChainError(err, "error saving vault index")
	}
	return nil
}

func (v Vault) ReadData(name string) ([]byte, error) {
	id, err := v.getID(name)
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(v.filepath(id))
	if err != nil {
		return nil, util.ChainError(err, "error reading file from vault")
	}

	return bytes, nil
}

func (v Vault) DeleteData(name string) error {
	id, err := v.getID(name)
	if err != nil {
		return err
	}

	err = os.Remove(v.filepath(id))
	if err != nil {
		return util.ChainError(err, "error deleting file from vault")
	}

	err = v.deleteIndex(name)
	if err != nil {
		return util.ChainError(err, "error deleting vault index")
	}
	return nil
}

func (v Vault) filepath(id uuid.UUID) string {
	return filepath.Join(v.Path, id.String()+".crypt")
}

func (v Vault) getID(name string) (uuid.UUID, error) {
	id, ok := v.index[name]
	if !ok {
		return uuid.Nil, util.Error(fmt.Sprintf("name \"%s\" not found", name))
	}
	return id, nil
}
