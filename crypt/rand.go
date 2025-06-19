package crypt

import (
	"crypto/rand"
)

func randSalt() []byte {
	salt := make([]byte, SALT_SIZE)
	rand.Read(salt)

	return salt
}

func randNonce() []byte {
	nonce := make([]byte, NONCE_SIZE)
	rand.Read(nonce)

	return nonce
}
