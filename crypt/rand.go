package crypt

import (
	"crypto/rand"

	"github.com/binary-soup/go-command/util"
)

func randSalt() ([]byte, error) {
	return randBytes("salt", SALT_SIZE)
}

func randNonce() ([]byte, error) {
	return randBytes("nonce", NONCE_SIZE)
}

func randBytes(name string, size int) ([]byte, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return nil, util.ChainErrorF(err, "error creating rand %s bytes", name)
	}

	return bytes, nil
}
