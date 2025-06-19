package crypt

import "github.com/binary-soup/go-command/util"

type Ciphertext []byte

const (
	NONCE_SIZE = 12 // AES-GCM standard nonce
)

func NewCiphertext(nonce, text []byte) Ciphertext {
	ciphertext := make(Ciphertext, NONCE_SIZE+len(text))
	copy(ciphertext.Nonce(), nonce)
	copy(ciphertext.Text(), text)

	return ciphertext
}

func LoadCiphertext(bytes []byte) (Header, error) {
	if len(bytes) <= NONCE_SIZE {
		return nil, util.Error("data shorter than ciphertext min size")
	}

	return bytes, nil
}

func (text Ciphertext) Nonce() []byte {
	return text[:NONCE_SIZE]
}

func (text Ciphertext) Text() []byte {
	return text[NONCE_SIZE:]
}
