package crypt

import (
	"crypto/rand"
	"io"

	"github.com/binary-soup/go-command/util"
)

func Encrypt(plaintext []byte) ([]byte, error) {
	gcm, err := NewAESGCM(KEY)
	if err != nil {
		return nil, err
	}

	// generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, util.ChainError(err, "error generating random nonce")
	}

	// encrypt plaintext
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return append(nonce, ciphertext...), nil
}
