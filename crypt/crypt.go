package crypt

import "github.com/binary-soup/go-command/util"

func (c Crypt) Encrypt(plaintext []byte) ([]byte, error) {
	nonce, err := randNonce()
	if err != nil {
		return nil, err
	}

	ciphertext := c.Cipher.Seal(nil, nonce, plaintext, nil)

	return buildCipherText(c.Salt, nonce, ciphertext), nil
}

func (c Crypt) Decrypt(bytes []byte) ([]byte, error) {
	nonceStart := SALT_SIZE + NONCE_SIZE

	if len(bytes) < nonceStart {
		return nil, util.Error("data too short for nonce")
	}
	nonce, ciphertext := bytes[SALT_SIZE:nonceStart], bytes[nonceStart:]

	plaintext, err := c.Cipher.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, util.ChainError(err, "error decrypting bytes")
	}

	return plaintext, nil
}

func buildCipherText(data ...[]byte) []byte {
	size := 0
	for _, bytes := range data {
		size += len(bytes)
	}

	ciphertext := make([]byte, size)
	idx := 0

	for _, bytes := range data {
		for _, b := range bytes {
			ciphertext[idx] = b
			idx++
		}
	}

	return ciphertext
}
