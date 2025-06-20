package data

import (
	"encoding/json"
	"pvault/crypt"

	"github.com/binary-soup/go-command/util"
	"github.com/google/uuid"
)

type Password struct {
	Name          string         `json:"name"`
	Password      string         `json:"password"`
	Username      string         `json:"username,omitempty"`
	URL           string         `json:"url,omitempty"`
	RecoveryCodes []string       `json:"recovery_codes,omitempty"`
	Cache         *PasswordCache `json:"cache,omitempty"`
}

type PasswordCache struct {
	Passkey string    `json:"passkey"`
	ID      uuid.UUID `json:"uuid"`
}

func NewPasswordCache(passkey string) *PasswordCache {
	return &PasswordCache{
		Passkey: passkey,
		ID:      uuid.New(),
	}
}

func LoadPasswordFile(path string) (*Password, error) {
	return util.LoadJSON[Password]("password", path)
}

func DecryptPassword(c *crypt.Crypt, ciphertext crypt.Ciphertext) (*Password, error) {
	plaintext, err := c.Decrypt(ciphertext)
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

func (password Password) Validate() error {
	if password.Name == "" {
		return util.Error("\"name\" cannot be empty")
	}
	if password.Password == "" && len(password.RecoveryCodes) == 0 {
		return util.Error("both \"password\" and \"recovery codes\" cannot be empty")
	}
	return nil
}

func (password Password) SaveToFile(path string) error {
	return util.SaveJSON("password", &password, path)
}

func (password Password) Encrypt(c *crypt.Crypt) (crypt.Ciphertext, error) {
	plaintext, err := json.Marshal(password)
	if err != nil {
		return nil, util.ChainError(err, "error marshaling password JSON")
	}

	return c.Encrypt(plaintext), nil
}
