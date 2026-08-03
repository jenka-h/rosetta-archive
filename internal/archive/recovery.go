package archive

import (
	"fmt"
	"io"
	"os"

	"rosetta-archive/internal/format"
)

// Recover validates an archive from disk. If footer/central-directory parsing
// fails, it scans local file records from the data region and returns recoverable
// entries whose local metadata is still structurally valid.
func Recover(archivePath string) ([]Entry, error) {
	if archivePath == "" {
		return nil, fmt.Errorf("recover archive: missing archive path")
	}
	f, err := os.Open(archivePath)
	if err != nil {
		return nil, fmt.Errorf("open archive: %w", err)
	}
	defer closeIfNeeded(f)
	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat archive: %w", err)
	}
	return RecoverReader(f, info.Size())
}

// RecoverReader validates an archive from an existing random-access reader. It
// falls back to sequential local-record scanning when the central directory or
// footer cannot be parsed.
func RecoverReader(r io.ReaderAt, size int64) ([]Entry, error) {
	reader, err := NewReader(r, size)
	if err == nil {
		return reader.Entries(), nil
	}
	entries, scanErr := scanLocalEntries(r, size)
	if scanErr != nil {
		return nil, fmt.Errorf("normal parse failed: %v; recovery scan failed: %w", err, scanErr)
	}
	return entries, nil
}

func scanLocalEntries(r io.ReaderAt, size int64) ([]Entry, error) {
	if size < int64(format.ArchiveHeaderSize) {
		return nil, fmt.Errorf("archive too small")
	}
	if _, err := format.DecodeHeader(io.NewSectionReader(r, 0, int64(format.ArchiveHeaderSize))); err != nil {
		return nil, fmt.Errorf("decode header: %w", err)
	}
	limit := uint64(size)
	if size > int64(format.FooterSize) {
		limit = uint64(size - int64(format.FooterSize))
	}
	offset := uint64(format.ArchiveHeaderSize)
	entries := make([]Entry, 0)
	seen := make(map[string]struct{})
	for offset+uint64(format.EntryHeaderSize) <= limit {
		entry, next, err := scanOneEntry(r, offset, limit)
		if err != nil {
			break
		}
		if _, ok := seen[entry.Path]; ok {
			break
		}
		seen[entry.Path] = struct{}{}
		entries = append(entries, entry)
		offset = next
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no recoverable entries")
	}
	return entries, nil
}

func scanOneEntry(r io.ReaderAt, offset uint64, limit uint64) (Entry, uint64, error) {
	section := io.NewSectionReader(r, int64(offset), int64(limit-offset))
	local, err := format.DecodeEntry(section)
	if err != nil {
		return Entry{}, offset, err
	}
	dataOffset := offset + uint64(format.EntryHeaderSize) + uint64(local.Header.PathLength) + uint64(local.Header.ExtraMetadataLength)
	end := dataOffset + local.Header.CompressedSize
	if end < dataOffset || end > limit {
		return Entry{}, offset, fmt.Errorf("entry payload outside recoverable range")
	}
	return Entry{
		Path:       local.Path,
		Type:       local.Header.EntryType,
		Size:       local.Header.UncompressedSize,
		CRC32:      local.Header.DataCRC32,
		Offset:     offset,
		DataOffset: dataOffset,
	}, end, nil
}
