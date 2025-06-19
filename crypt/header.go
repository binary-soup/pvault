package crypt

const (
	HASH_SIZE   = 60
	SALT_SIZE   = 16
	HEADER_SIZE = HASH_SIZE + SALT_SIZE
)

type Header []byte

func EmptyHeader() Header {
	return make(Header, HEADER_SIZE)
}

func NewHeader(hash, salt []byte) Header {
	header := EmptyHeader()
	copy(header.Hash(), hash)
	copy(header.Salt(), salt)

	return header
}

func (header Header) Hash() []byte {
	return header[:HASH_SIZE]
}

func (header Header) Salt() []byte {
	return header[HASH_SIZE:]
}
