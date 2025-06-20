package vault

import (
	"io"
	"os"
	"path/filepath"
	"pvault/crypt"

	"github.com/binary-soup/go-command/util"
	"github.com/google/uuid"
)

func (v Vault) SaveData(header crypt.Header, ciphertext crypt.Ciphertext, id uuid.UUID, name string) error {
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

	v.Index.AddPair(name, id)
	return nil
}

func (v Vault) ReadData(name string) ([]byte, []byte, uuid.UUID, error) {
	id, err := v.Index.GetID(name)
	if err != nil {
		return nil, nil, uuid.Nil, err
	}

	file, err := os.Open(v.filepath(id))
	if err != nil {
		return nil, nil, id, util.ChainError(err, "error opening vault file")
	}
	defer file.Close()

	header := crypt.EmptyHeader()

	_, err = file.Read(header)
	if err != nil {
		return nil, nil, id, util.ChainError(err, "error reading header from vault")
	}

	ciphertext, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, id, util.ChainError(err, "error reading ciphertext from vault")
	}

	return header, ciphertext, id, nil
}

func (v Vault) DeleteData(name string) error {
	id, err := v.Index.GetID(name)
	if err != nil {
		return err
	}

	err = os.Remove(v.filepath(id))
	if err != nil {
		return util.ChainError(err, "error deleting file from vault")
	}

	v.Index.DeleteName(name)
	return nil
}

func (v Vault) filepath(id uuid.UUID) string {
	return filepath.Join(v.Path, id.String()+".crypt")
}
