package vault

import (
	"io"
	"os"
	"path/filepath"
	"pvault/crypt"

	"github.com/binary-soup/go-command/alert"
	"github.com/google/uuid"
)

func (v Vault) SaveRaw(id uuid.UUID, name string, bytes []byte) error {
	err := os.WriteFile(v.filepath(id), bytes, 0666)
	if err != nil {
		return alert.ChainError(err, "error writing raw vault file")
	}

	v.Index.AddPair(name, id)
	return nil
}

func (v Vault) SaveData(header crypt.Header, ciphertext crypt.Ciphertext, id uuid.UUID, name string) error {
	file, err := os.Create(v.filepath(id))
	if err != nil {
		return alert.ChainError(err, "error creating vault file")
	}
	defer file.Close()

	_, err = file.Write(header)
	if err != nil {
		return alert.ChainError(err, "error writing header to vault")
	}

	_, err = file.Write(ciphertext)
	if err != nil {
		return alert.ChainError(err, "error writing ciphertext to vault")
	}

	v.Index.AddPair(name, id)
	return nil
}

func (v Vault) LoadRaw(id uuid.UUID) ([]byte, error) {
	raw, err := os.ReadFile(v.filepath(id))
	if err != nil {
		return nil, alert.ChainError(err, "error reading raw vault file")
	}
	return raw, nil
}

func (v Vault) LoadData(id uuid.UUID) ([]byte, []byte, error) {
	file, err := os.Open(v.filepath(id))
	if err != nil {
		return nil, nil, alert.ChainError(err, "error opening vault file")
	}
	defer file.Close()

	header := crypt.EmptyHeader()

	_, err = file.Read(header)
	if err != nil {
		return nil, nil, alert.ChainError(err, "error reading header from vault")
	}

	ciphertext, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, alert.ChainError(err, "error reading ciphertext from vault")
	}

	return header, ciphertext, nil
}

func (v Vault) DeleteData(id uuid.UUID) error {
	err := os.Remove(v.filepath(id))
	if err != nil {
		return alert.ChainError(err, "error deleting file from vault")
	}

	v.Index.DeleteID(id)
	return nil
}

func (v Vault) filepath(id uuid.UUID) string {
	return filepath.Join(v.Path, id.String()+".crypt")
}
