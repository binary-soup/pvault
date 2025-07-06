package bytez

import (
	"encoding/binary"
	"io"

	"github.com/binary-soup/go-command/alert"
)

type Writable interface {
	ToBytes() []byte
}

func writeHeader(name string, dest io.Writer, count uint32) error {
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, count)

	_, err := dest.Write(header)
	if err != nil {
		return alert.ChainErrorF(err, "error writing %s header", name)
	}
	return nil
}

func WriteBlobRaw(dest io.Writer, blob []byte) error {
	err := writeHeader("blob", dest, uint32(len(blob)))
	if err != nil {
		return err
	}

	_, err = dest.Write(blob)
	if err != nil {
		return alert.ChainError(err, "error writing blob")
	}

	return nil
}

func WriteBlob(dest io.Writer, w Writable) error {
	return WriteBlobRaw(dest, w.ToBytes())
}

func WriteBlockFunc(dest io.Writer, count int, blob func(int) []byte) error {
	err := writeHeader("block", dest, uint32(count))
	if err != nil {
		return err
	}

	for i := range count {
		err := WriteBlobRaw(dest, blob(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteBlockRaw(dest io.Writer, block [][]byte) error {
	return WriteBlockFunc(dest, len(block), func(i int) []byte {
		return block[i]
	})
}

func WriteBlock(dest io.Writer, block []Writable) error {
	return WriteBlockFunc(dest, len(block), func(i int) []byte {
		return block[i].ToBytes()
	})
}
