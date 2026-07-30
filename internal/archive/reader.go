package archive

import "io"

// Reader owns parsed archive state and provides shared validation for all commands.
type Reader struct {
	entries []Entry
}

// Open parses and validates a ROSA archive from disk.
//
// TODO: implement header, metadata, payload range, footer, and checksum validation.
func Open(archivePath string) (*Reader, error) {
	if archivePath == "" {
		return nil, notImplemented("open archive: missing archive path")
	}
	return nil, notImplemented("open archive")
}

// NewReader parses and validates a ROSA archive from an existing random-access reader.
//
// TODO: use io.ReaderAt/io.Seeker based parsing once the binary layout is defined.
func NewReader(r io.ReaderAt, size int64) (*Reader, error) {
	if r == nil {
		return nil, notImplemented("open archive reader: nil reader")
	}
	if size < 0 {
		return nil, notImplemented("open archive reader: negative size")
	}
	return nil, notImplemented("open archive reader")
}

// Entries returns the parsed archive entries in archive order.
func (r *Reader) Entries() []Entry {
	if r == nil || len(r.entries) == 0 {
		return nil
	}
	entries := make([]Entry, len(r.entries))
	copy(entries, r.entries)
	return entries
}
