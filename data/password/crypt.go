package password

import (
	"encoding/json"
	"pvault/crypt"

	"github.com/binary-soup/go-commando/alert"
)

func (password Password) Encrypt(c *crypt.Crypt) (crypt.Ciphertext, error) {
	plaintext, err := json.Marshal(password)
	if err != nil {
		return nil, alert.ChainError(err, "error marshaling password JSON")
	}

	return c.Encrypt(plaintext), nil
}

func Decrypt(c *crypt.Crypt, ciphertext crypt.Ciphertext) (Password, error) {
	var password Password

	plaintext, err := c.Decrypt(ciphertext)
	if err != nil {
		return password, alert.ChainError(err, "error decrypting password")
	}

	err = json.Unmarshal(plaintext, &password)
	if err != nil {
		return password, alert.ChainError(err, "error unmarshaling password JSON")
	}

	return password, nil
}
