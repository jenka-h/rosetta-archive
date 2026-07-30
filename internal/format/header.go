package format

import "io"

type Header struct{}

// TODO: define the header fields, little-endian encoding, reserved fields, and CRC32 coverage.
func EncodeHeader(w io.Writer, h Header) error {
	if w == nil {
		return ErrNilWriter
	}
	return ErrFormatNotDefined
}

// TODO: validate magic, version, feature flags, size, reserved fields, and header CRC32.
func DecodeHeader(r io.Reader) (Header, error) {
	if r == nil {
		return Header{}, ErrNilReader
	}
	return Header{}, ErrFormatNotDefined
}
