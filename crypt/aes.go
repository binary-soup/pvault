package crypt

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/binary-soup/go-command/util"
)

func NewAESGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, util.ChainError(err, "error creating AES cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, util.ChainError(err, "error creating GCM mode")
	}

	return gcm, nil
}
