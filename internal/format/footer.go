package format

import "io"

// Footer is a placeholder for the future fixed-size archive footer.
type Footer struct{}

// EncodeFooter writes the archive footer.
//
// TODO: define footer fields, metadata checksum, archive checksum, footer CRC32, and location rules.
func EncodeFooter(w io.Writer, f Footer) error {
	if w == nil {
		return ErrNilWriter
	}
	return ErrFormatNotDefined
}

// DecodeFooter reads and validates the archive footer.
//
// TODO: validate footer location, version compatibility, reserved fields, and checksum coverage.
func DecodeFooter(r io.Reader) (Footer, error) {
	if r == nil {
		return Footer{}, ErrNilReader
	}
	return Footer{}, ErrFormatNotDefined
}
