package crypt

import (
	"crypto/rand"
	"math/big"

	"github.com/binary-soup/go-command/alert"
)

const CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}<>?"

func RandPassword(length int) (string, error) {
	if length < 1 {
		return "", alert.Error("length too short")
	}

	bytes := make([]byte, length)
	max := big.NewInt(int64(len(CHARSET)))

	for i := range bytes {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", alert.ChainError(err, "error generating rand int")
		}

		bytes[i] = CHARSET[num.Int64()]
	}

	return string(bytes), nil
}
