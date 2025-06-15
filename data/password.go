package data

import (
	"encoding/json"
	"os"
	"passwords/crypt"

	"github.com/binary-soup/go-command/util"
)

type Password struct {
	Username      string   `json:"username"`
	Password      string   `json:"password"`
	RecoveryCodes []string `json:"recovery_codes"`
}

func LoadPasswordFile(path string) (*Password, error) {
	return util.LoadJSON[Password]("password", path)
}

func DecryptPasswordFromFile(path string) (*Password, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, util.ChainError(err, "error reading crypt file")
	}

	plaintext, err := crypt.Decrypt(data)
	if err != nil {
		return nil, util.ChainError(err, "error decrypting password")
	}

	password := &Password{}

	err = json.Unmarshal(plaintext, password)
	if err != nil {
		return nil, util.ChainError(err, "error unmarshaling password JSON")
	}

	return password, nil
}

func (password Password) SaveToFile(path string) error {
	return util.SaveJSON("password", &password, path)
}

func (password Password) EncryptToFile(path string) error {
	plaintext, err := json.Marshal(password)
	if err != nil {
		return util.ChainError(err, "error marshaling password JSON")
	}

	ciphertext, err := crypt.Encrypt(plaintext)
	if err != nil {
		return util.ChainError(err, "error encrypting password")
	}

	err = os.WriteFile(path, ciphertext, 0600)
	if err != nil {
		return util.ChainError(err, "error saving crypt file")
	}

	return nil
}
