package bytez

import (
	"encoding/binary"
	"io"

	"github.com/binary-soup/go-commando/alert"
)

type Readable[T any] interface {
	FromBytes(bytes []byte) (T, error)
}

func readHeader(name string, src io.Reader) (uint32, error) {
	header := make([]byte, 4)

	_, err := io.ReadFull(src, header)
	if err != nil {
		return 0, alert.ChainErrorF(err, "error reading %s header", name)
	}

	return binary.BigEndian.Uint32(header), nil
}

func ReadBlobRaw(src io.Reader) ([]byte, error) {
	count, err := readHeader("blob", src)
	if err != nil {
		return nil, err
	}
	blob := make([]byte, count)

	_, err = io.ReadFull(src, blob)
	if err != nil {
		return nil, alert.ChainError(err, "error reading blob")
	}

	return blob, nil
}

func ReadBlob[T any](src io.Reader, r Readable[T]) (T, error) {
	var blob T

	bytes, err := ReadBlobRaw(src)
	if err != nil {
		return blob, err
	}

	blob, err = r.FromBytes(bytes)
	if err != nil {
		return blob, alert.ChainError(err, "error parsing blob")
	}

	return blob, nil
}

func ReadBlockFunc[T any](src io.Reader, parse func([]byte) (T, error)) ([]T, error) {
	count, err := readHeader("block", src)
	if err != nil {
		return nil, err
	}
	block := make([]T, count)

	for i := range count {
		blob, err := ReadBlobRaw(src)
		if err != nil {
			return nil, alert.ChainErrorF(err, "error reading blob at index %d", i)
		}

		block[i], err = parse(blob)
		if err != nil {
			return nil, alert.ChainErrorF(err, "error parsing blob at index %d", i)
		}
	}

	return block, nil
}

func ReadBlockRaw(src io.Reader) ([][]byte, error) {
	return ReadBlockFunc(src, func(blob []byte) ([]byte, error) {
		return blob, nil
	})
}

func ReadBlock[T any](src io.Reader, r Readable[T]) ([]T, error) {
	return ReadBlockFunc(src, func(blob []byte) (T, error) {
		return r.FromBytes(blob)
	})
}
