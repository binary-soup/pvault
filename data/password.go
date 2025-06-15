package data

import (
	"encoding/json"
	"passwords/crypt"

	"github.com/binary-soup/go-command/util"
)

type Password struct {
	Username      string   `json:"username"`
	Password      string   `json:"password"`
	URL           string   `json:"url,omitempty"`
	RecoveryCodes []string `json:"recovery_codes,omitempty"`
}

func LoadPasswordFile(path string) (*Password, error) {
	return util.LoadJSON[Password]("password", path)
}

func DecryptPassword(bytes []byte) (*Password, error) {
	plaintext, err := crypt.Decrypt(bytes)
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

func (password Password) Encrypt() ([]byte, error) {
	plaintext, err := json.Marshal(password)
	if err != nil {
		return nil, util.ChainError(err, "error marshaling password JSON")
	}

	ciphertext, err := crypt.Encrypt(plaintext)
	if err != nil {
		return nil, util.ChainError(err, "error encrypting password")
	}

	return ciphertext, nil
}
