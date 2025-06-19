package crypt

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

func (text Ciphertext) Nonce() []byte {
	return text[:NONCE_SIZE]
}

func (text Ciphertext) Text() []byte {
	return text[NONCE_SIZE:]
}
