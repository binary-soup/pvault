package crypt

import (
	"github.com/binary-soup/go-command/util"
)

func Decrypt(bytes []byte) ([]byte, error) {
	gcm, err := NewAESGCM(KEY)
	if err != nil {
		return nil, err
	}

	// extract nonce
	nonceSize := gcm.NonceSize()
	if len(bytes) < nonceSize {
		return nil, util.ChainError(err, "data too short")
	}
	nonce, ciphertext := bytes[:nonceSize], bytes[nonceSize:]

	// decrypt ciphertext
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, util.ChainError(err, "error decrypting bytes")
	}

	return plaintext, nil
}
