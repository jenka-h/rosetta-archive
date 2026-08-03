package archive

import (
	"bytes"
	"fmt"
	"io"

	"rosetta-archive/internal/format"
)

// Stream visits entries from an io.Reader. The current ROSA layout stores the
// central directory location in the footer, so this function buffers the input in
// order to reuse the same validated parser as Open/NewReader.
func Stream(src io.Reader, visit func(Entry, io.Reader) error) error {
	if src == nil {
		return fmt.Errorf("stream archive: nil reader")
	}
	if visit == nil {
		return fmt.Errorf("stream archive: nil visit callback")
	}
	data, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("read archive stream: %w", err)
	}
	r, err := NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	for _, meta := range r.directory {
		entry := entryFromDirectory(meta)
		var payload io.Reader
		if meta.EntryType == format.EntryTypeFile {
			decoded, err := r.readPayload(meta)
			if err != nil {
				return err
			}
			payload = bytes.NewReader(decoded)
		}
		if err := visit(entry, payload); err != nil {
			return err
		}
	}
	return nil
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
	meta := r.directory[i]
	entry := entryFromDirectory(meta)
	if meta.EntryType == format.EntryTypeDirectory {
		return nil, entry, fmt.Errorf("open archive entry: %q is a directory", archivePath)
	}
	payload, err := r.readPayload(meta)
	if err != nil {
		return nil, Entry{}, err
	}
	return bytes.NewReader(payload), entry, nil
}
