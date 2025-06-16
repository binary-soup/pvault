package crypt

import "github.com/binary-soup/go-command/util"

type DataBlock []byte

func BuildDataBlock(hash, salt, nonce, ciphertext []byte) DataBlock {
	block := make(DataBlock, len(hash)+len(salt)+len(nonce)+len(ciphertext))
	idx := 0

	for _, bytes := range [][]byte{hash, salt, nonce, ciphertext} {
		for _, b := range bytes {
			block[idx] = b
			idx++
		}
	}

	return block
}

func LoadDataBlock(bytes []byte) (DataBlock, error) {
	if len(bytes) < BCRYPT_HASH_SIZE+SALT_SIZE+NONCE_SIZE {
		return nil, util.Error("data block invalid")
	}

	return bytes, nil
}

func (block DataBlock) BCryptHash() []byte {
	return block[:BCRYPT_HASH_SIZE]
}

func (block DataBlock) Salt() []byte {
	return block[BCRYPT_HASH_SIZE : BCRYPT_HASH_SIZE+SALT_SIZE]
}

func (block DataBlock) Nonce() []byte {
	return block[BCRYPT_HASH_SIZE+SALT_SIZE : BCRYPT_HASH_SIZE+SALT_SIZE+NONCE_SIZE]
}

func (block DataBlock) Ciphertext() []byte {
	return block[BCRYPT_HASH_SIZE+SALT_SIZE+NONCE_SIZE:]
}
