package crypt

import "github.com/binary-soup/go-command/util"

func (c Crypt) Encrypt(plaintext []byte) Ciphertext {
	nonce := randNonce()
	text := c.Cipher.Seal(nil, nonce, plaintext, nil)

	return NewCiphertext(nonce, text)
}

func (c Crypt) Decrypt(ciphertext Ciphertext) ([]byte, error) {
	plaintext, err := c.Cipher.Open(nil, ciphertext.Nonce(), ciphertext.Text(), nil)
	if err != nil {
		return nil, util.ChainError(err, "error decrypting bytes")
	}

	return plaintext, nil
}
