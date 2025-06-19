package vault

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
	"passwords/crypt"
	"path/filepath"

	"github.com/binary-soup/go-command/util"
)

func (v Vault) SaveData(header crypt.Header, ciphertext crypt.Ciphertext, name string) error {
	file, err := os.Create(v.filepath(name))
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

	v.index.Add(name)
	return nil
}

func (v Vault) ReadData(name string) ([]byte, []byte, error) {
	file, err := os.Open(v.filepath(name))
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
	err := os.Remove(v.filepath(name))
	if err != nil {
		return util.ChainError(err, "error deleting file from vault")
	}

	v.index.Delete(name)
	return nil
}

func (v Vault) filepath(name string) string {
	hash := sha256.Sum256([]byte(name))
	return filepath.Join(v.Path, base64.RawURLEncoding.EncodeToString(hash[:])+".crypt")
}
