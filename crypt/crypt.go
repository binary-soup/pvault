package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/sha256"

	"github.com/binary-soup/go-command/alert"
	"golang.org/x/crypto/bcrypt"
)

const (
	KEY_SIZE          = 32 // AES-256
	PBKDF2_ITERATIONS = 100_000
)

type Crypt struct {
	Header Header
	cipher cipher.AEAD
}

func NewCrypt(passkey string) (*Crypt, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passkey), bcrypt.DefaultCost)
	if err != nil {
		return nil, alert.ChainError(err, "error hashing passkey")
	}

	return newCrypt(passkey, NewHeader(hash, randSalt()))
}

func LoadCrypt(passkey string, header Header) (*Crypt, bool, error) {
	err := bcrypt.CompareHashAndPassword(header.Hash(), []byte(passkey))
	if err != nil {
		return nil, true, nil
	}

	c, err := newCrypt(passkey, header)
	return c, false, err
}

func newCrypt(passkey string, header Header) (*Crypt, error) {
	key, err := pbkdf2.Key(sha256.New, passkey, header.Salt(), PBKDF2_ITERATIONS, KEY_SIZE)
	if err != nil {
		return nil, alert.ChainError(err, "error generating key from passkey")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, alert.ChainError(err, "error creating AES cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, alert.ChainError(err, "error creating GCM mode")
	}

	return &Crypt{
		Header: header,
		cipher: gcm,
	}, nil
}
