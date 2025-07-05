package crypt

import "github.com/binary-soup/go-command/alert"

func (c Crypt) Encrypt(plaintext []byte) Ciphertext {
	nonce := randNonce()
	text := c.cipher.Seal(nil, nonce, plaintext, nil)

	return NewCiphertext(nonce, text)
}

func (c Crypt) Decrypt(ciphertext Ciphertext) ([]byte, error) {
	plaintext, err := c.cipher.Open(nil, ciphertext.Nonce(), ciphertext.Text(), nil)
	if err != nil {
		return nil, alert.ChainError(err, "error decrypting bytes")
	}

	return plaintext, nil
}
