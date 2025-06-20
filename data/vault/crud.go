package vault

import (
	"fmt"
	"io"
	"os"
	"passwords/crypt"
	"path/filepath"

	"github.com/binary-soup/go-command/util"
	"github.com/google/uuid"
)

func (v Vault) CreateData(header crypt.Header, ciphertext crypt.Ciphertext, name string) error {
	id := uuid.New()

	file, err := os.Create(v.filepath(id))
	if err != nil {
		return util.ChainError(err, "error creating vault file")
	}
	defer file.Close()

	_, err = file.Write(header)
	if err != nil {
		return util.ChainError(err, "error writing header to vault")
	}

	_, err = file.Write(ciphertext)
	if err != nil {
		return util.ChainError(err, "error writing ciphertext to vault")
	}

	v.index[name] = id
	return nil
}

func (v Vault) ReadData(name string) ([]byte, []byte, error) {
	id, err := v.getID(name)
	if err != nil {
		return nil, nil, err
	}

	file, err := os.Open(v.filepath(id))
	if err != nil {
		return nil, nil, util.ChainError(err, "error opening vault file")
	}
	defer file.Close()

	header := crypt.EmptyHeader()

	_, err = file.Read(header)
	if err != nil {
		return nil, nil, util.ChainError(err, "error reading header from vault")
	}

	ciphertext, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, util.ChainError(err, "error reading ciphertext from vault")
	}

	return header, ciphertext, nil
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

	delete(v.index, name)
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
