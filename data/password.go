package data

import (
	"encoding/json"
	"passwords/crypt"

	"github.com/binary-soup/go-command/util"
)

type Password struct {
	Name          string   `json:"name"`
	Password      string   `json:"password"`
	Username      string   `json:"username,omitempty"`
	URL           string   `json:"url,omitempty"`
	RecoveryCodes []string `json:"recovery_codes,omitempty"`
}

func LoadPasswordFile(path string) (*Password, error) {
	return util.LoadJSON[Password]("password", path)
}

func DecryptPassword(c *crypt.Crypt, bytes []byte) (*Password, error) {
	plaintext, err := c.Decrypt(bytes)
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

func (password Password) Encrypt(c *crypt.Crypt) ([]byte, error) {
	plaintext, err := json.Marshal(password)
	if err != nil {
		return nil, util.ChainError(err, "error marshaling password JSON")
	}

	ciphertext, err := c.Encrypt(plaintext)
	if err != nil {
		return nil, util.ChainError(err, "error encrypting password")
	}

	return ciphertext, nil
}
