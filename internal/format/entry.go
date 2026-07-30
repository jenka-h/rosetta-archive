package format

import "io"

// Entry is a placeholder for a future metadata record.
type Entry struct{}

// EncodeEntry writes one metadata entry.
//
// TODO: define record type, path encoding, offsets, sizes, CRC32, and reserved fields.
func EncodeEntry(w io.Writer, e Entry) error {
	if w == nil {
		return ErrNilWriter
	}
	return ErrFormatNotDefined
}

// DecodeEntry reads and validates one metadata entry.
//
// TODO: validate record length, path length, UTF-8, entry type, offsets, sizes, and CRC32.
func DecodeEntry(r io.Reader) (Entry, error) {
	if r == nil {
		return Entry{}, ErrNilReader
	}
	return Entry{}, ErrFormatNotDefined
}
