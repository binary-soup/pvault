package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/sha256"

	"github.com/binary-soup/go-command/util"
	"golang.org/x/crypto/bcrypt"
)

const (
	BCRYPT_HASH_SIZE  = 60
	SALT_SIZE         = 16
	NONCE_SIZE        = 12 // AES-GCM standard nonce
	KEY_SIZE          = 32 // AES-256
	PBKDF2_ITERATIONS = 100_000
)

type Crypt struct {
	PasskeyHash []byte
	Salt        []byte
	Cipher      cipher.AEAD
}

func NewCrypt(passkey string) (*Crypt, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passkey), bcrypt.DefaultCost)
	if err != nil {
		return nil, util.ChainError(err, "error hashing passkey")
	}

	salt, err := randSalt()
	if err != nil {
		return nil, err
	}

	return newCrypt(passkey, hash, salt)
}

func LoadCrypt(passkey string, bytes []byte) (*Crypt, bool, error) {
	block, err := LoadDataBlock(bytes)
	if err != nil {
		return nil, false, err
	}

	err = bcrypt.CompareHashAndPassword(block.BCryptHash(), []byte(passkey))
	if err != nil {
		return nil, true, nil
	}

	c, err := newCrypt(passkey, block.BCryptHash(), block.Salt())
	return c, false, err
}

func newCrypt(passkey string, hash, salt []byte) (*Crypt, error) {
	key, err := pbkdf2.Key(sha256.New, passkey, salt, PBKDF2_ITERATIONS, KEY_SIZE)
	if err != nil {
		return nil, util.ChainError(err, "error generating key from passkey")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, util.ChainError(err, "error creating AES cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, util.ChainError(err, "error creating GCM mode")
	}

	return &Crypt{
		PasskeyHash: hash,
		Salt:        salt,
		Cipher:      gcm,
	}, nil
}
