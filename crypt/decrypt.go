package crypt

import (
	"github.com/binary-soup/go-command/util"
)

func Decrypt(data []byte) ([]byte, error) {
	gcm, err := NewAESGCM(KEY)
	if err != nil {
		return nil, err
	}

	// extract nonce
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, util.ChainError(err, "data too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// decrypt ciphertext
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, util.ChainError(err, "error decrypting data")
	}

	return plaintext, nil
}
