package tcp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	BinaryType uint8 = iota + 1
	StringType

	// MaxPayloadSize 10 MB
	MaxPayloadSize uint32 = 10 << 20
)

var ErrMaxPayloadSize = errors.New("max payload size exceeded")

type Payload interface {
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes() []byte
}

type Binary []byte

func (m Binary) Bytes() []byte {
	return m
}

func (m Binary) String() string {
	return string(m)
}

func (m Binary) WriteTo(w io.Writer) (int64, error) {
	if err := binary.Write(w, binary.BigEndian, BinaryType); err != nil {
		return 0, err
	}
	var n int64 = 1
	if err := binary.Write(w, binary.BigEndian, uint32(len(m))); err != nil {
		return n, err
	}
	n += 4
	o, err := w.Write(m)

	return n + int64(o), err
}
