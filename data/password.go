package data

import (
	"encoding/json"
	"passwords/crypt"

	"github.com/binary-soup/go-command/util"
)

const PASSKEY = "Password123!"

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
	c, err := crypt.LoadCrypt(PASSKEY, bytes)
	if err != nil {
		return nil, err
	}

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

func (password Password) Encrypt() ([]byte, error) {
	plaintext, err := json.Marshal(password)
	if err != nil {
		return nil, util.ChainError(err, "error marshaling password JSON")
	}

	c, err := crypt.NewCrypt(PASSKEY)
	if err != nil {
		return nil, util.ChainError(err, "error initializing crypt tool")
	}

	ciphertext, err := c.Encrypt(plaintext)
	if err != nil {
		return nil, util.ChainError(err, "error encrypting password")
	}

	return ciphertext, nil
}
