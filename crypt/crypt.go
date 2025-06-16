package crypt

import "github.com/binary-soup/go-command/util"

func (c Crypt) Encrypt(plaintext []byte) (DataBlock, error) {
	nonce, err := randNonce()
	if err != nil {
		return nil, err
	}

	ciphertext := c.Cipher.Seal(nil, nonce, plaintext, nil)

	return BuildDataBlock(c.PasskeyHash, c.Salt, nonce, ciphertext), nil
}

func (c Crypt) Decrypt(bytes []byte) ([]byte, error) {
	block, err := LoadDataBlock(bytes)
	if err != nil {
		return nil, err
	}

	plaintext, err := c.Cipher.Open(nil, block.Nonce(), block.Ciphertext(), nil)
	if err != nil {
		return nil, util.ChainError(err, "error decrypting bytes")
	}

	return plaintext, nil
}
