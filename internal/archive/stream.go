package archive

import (
	"fmt"
	"io"

	"rosetta-archive/internal/format"
)

type StreamOptions struct {
	VerifyDataCRC bool
}

type StreamEntry struct {
	Entry Entry
	Data  io.Reader
}

func Stream(r io.Reader, options StreamOptions, visit func(StreamEntry) error) error {
	if r == nil {
		return fmt.Errorf("stream archive: nil reader")
	}
	if visit == nil {
		return fmt.Errorf("stream archive: nil visit callback")
	}
	return notImplemented("stream archive")
}

func (r *Reader) OpenEntry(archivePath string) (io.Reader, Entry, error) {
	if r == nil {
		return nil, Entry{}, fmt.Errorf("open archive entry: nil reader")
	}
	if archivePath == "" {
		return nil, Entry{}, fmt.Errorf("open archive entry: missing path")
	}
	i, ok := r.index[archivePath]
	if !ok {
		return nil, Entry{}, fmt.Errorf("open archive entry: not found %q", archivePath)
	}
	entry := r.directory[i]
	public := entryFromDirectory(entry)
	if entry.EntryType == format.EntryTypeDirectory {
		return nil, public, fmt.Errorf("open archive entry: %q is a directory", archivePath)
	}
	return io.NewSectionReader(r.r, int64(entry.DataOffset), int64(entry.CompressedSize)), public, nil
}
