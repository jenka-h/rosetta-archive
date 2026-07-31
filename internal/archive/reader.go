package archive

import (
	"bytes"
	"fmt"
	"io"
	"os"

	core "rosetta-archive/internal"
	"rosetta-archive/internal/format"
)

type Reader struct {
	r         io.ReaderAt
	closer    io.Closer
	size      int64
	header    format.ArchiveHeader
	footer    format.ArchiveFooter
	directory []format.DirectoryEntry
	index     map[string]int
}

func Open(archivePath string) (*Reader, error) {
	if archivePath == "" {
		return nil, fmt.Errorf("open archive: missing archive path")
	}
	f, err := os.Open(archivePath)
	if err != nil {
		return nil, fmt.Errorf("open archive: %w", err)
	}
	info, err := f.Stat()
	if err != nil {
		closeIfNeeded(f)
		return nil, fmt.Errorf("stat archive: %w", err)
	}
	reader, err := newReader(f, f, info.Size())
	if err != nil {
		closeIfNeeded(f)
		return nil, err
	}
	return reader, nil
}

func NewReader(r io.ReaderAt, size int64) (*Reader, error) {
	return newReader(r, nil, size)
}

func newReader(r io.ReaderAt, closer io.Closer, size int64) (*Reader, error) {
	if r == nil {
		return nil, fmt.Errorf("open archive reader: nil reader")
	}
	if size < int64(format.ArchiveHeaderSize)+int64(format.FooterSize) {
		return nil, fmt.Errorf("open archive reader: archive too small")
	}
	header, err := readHeader(r)
	if err != nil {
		return nil, err
	}
	footer, err := readFooter(r, size)
	if err != nil {
		return nil, err
	}
	if err := validateFooterRange(footer, size); err != nil {
		return nil, err
	}
	directory, err := readDirectory(r, footer)
	if err != nil {
		return nil, err
	}
	index := make(map[string]int, len(directory))
	for i, entry := range directory {
		if _, ok := index[entry.Path]; ok {
			return nil, fmt.Errorf("duplicate archive path %q", entry.Path)
		}
		if err := validateEntryRange(entry, uint64(size), footer.CentralDirectoryOffset); err != nil {
			return nil, fmt.Errorf("validate directory entry %q: %w", entry.Path, err)
		}
		if err := validateLocalEntry(r, entry); err != nil {
			return nil, fmt.Errorf("validate local entry %q: %w", entry.Path, err)
		}
		index[entry.Path] = i
	}
	return &Reader{r: r, closer: closer, size: size, header: header, footer: footer, directory: directory, index: index}, nil
}

func (r *Reader) Close() error {
	if r == nil || r.closer == nil {
		return nil
	}
	return r.closer.Close()
}

func (r *Reader) Entries() []Entry {
	if r == nil || len(r.directory) == 0 {
		return nil
	}
	entries := make([]Entry, len(r.directory))
	for i, entry := range r.directory {
		entries[i] = entryFromDirectory(entry)
	}
	return entries
}

func readHeader(r io.ReaderAt) (format.ArchiveHeader, error) {
	header, err := format.DecodeHeader(io.NewSectionReader(r, 0, int64(format.ArchiveHeaderSize)))
	if err != nil {
		return format.ArchiveHeader{}, fmt.Errorf("decode archive header: %w", err)
	}
	return header, nil
}

func readFooter(r io.ReaderAt, size int64) (format.ArchiveFooter, error) {
	footer, err := format.DecodeFooter(io.NewSectionReader(r, size-int64(format.FooterSize), int64(format.FooterSize)))
	if err != nil {
		return format.ArchiveFooter{}, fmt.Errorf("decode archive footer: %w", err)
	}
	return footer, nil
}

func validateFooterRange(footer format.ArchiveFooter, size int64) error {
	footerOffset := uint64(size) - uint64(format.FooterSize)
	if footer.EntryCount > core.MaxEntries {
		return fmt.Errorf("entry count %d exceeds maximum %d", footer.EntryCount, core.MaxEntries)
	}
	if footer.CentralDirectoryOffset < uint64(format.ArchiveHeaderSize) {
		return fmt.Errorf("central directory offset before payload region")
	}
	if footer.CentralDirectoryOffset > uint64(size) || footer.CentralDirectorySize > uint64(size)-footer.CentralDirectoryOffset {
		return fmt.Errorf("central directory range is outside archive")
	}
	if footer.CentralDirectoryOffset+footer.CentralDirectorySize != footerOffset {
		return fmt.Errorf("central directory does not end at footer")
	}
	return nil
}

func readDirectory(r io.ReaderAt, footer format.ArchiveFooter) ([]format.DirectoryEntry, error) {
	data := make([]byte, footer.CentralDirectorySize)
	if _, err := r.ReadAt(data, int64(footer.CentralDirectoryOffset)); err != nil {
		return nil, fmt.Errorf("read central directory: %w", err)
	}
	if len(data) < int(format.DirectoryHeaderSize) {
		return nil, fmt.Errorf("central directory too small")
	}
	header, err := format.DecodeDirectoryHeader(bytes.NewReader(data[:format.DirectoryHeaderSize]))
	if err != nil {
		return nil, err
	}
	if header.EntryCount != footer.EntryCount {
		return nil, fmt.Errorf("directory entry count mismatch")
	}
	if format.DirectoryCRC32(data) != header.DirectoryCRC32 {
		return nil, core.ErrInvalidDirectoryCRC32
	}
	br := bytes.NewReader(data[format.DirectoryHeaderSize:])
	entries := make([]format.DirectoryEntry, 0, header.EntryCount)
	for i := uint32(0); i < header.EntryCount; i++ {
		entry, err := format.DecodeDirectoryEntry(br)
		if err != nil {
			return nil, fmt.Errorf("decode directory entry %d: %w", i, err)
		}
		entries = append(entries, entry)
	}
	if br.Len() != 0 {
		return nil, fmt.Errorf("central directory has trailing bytes")
	}
	return entries, nil
}

func validateEntryRange(entry format.DirectoryEntry, archiveSize uint64, directoryOffset uint64) error {
	if entry.FileRecordOffset < uint64(format.ArchiveHeaderSize) || entry.FileRecordOffset >= directoryOffset {
		return fmt.Errorf("invalid file record offset")
	}
	if entry.DataOffset < entry.FileRecordOffset || entry.DataOffset > directoryOffset {
		return fmt.Errorf("invalid data offset")
	}
	if entry.CompressedSize > archiveSize-entry.DataOffset {
		return fmt.Errorf("payload range outside archive")
	}
	if entry.DataOffset+entry.CompressedSize > directoryOffset {
		return fmt.Errorf("payload overlaps central directory")
	}
	return nil
}

func validateLocalEntry(r io.ReaderAt, dirEntry format.DirectoryEntry) error {
	localReader := io.NewSectionReader(r, int64(dirEntry.FileRecordOffset), int64(format.EntryHeaderSize)+int64(dirEntry.PathLength))
	local, err := format.DecodeEntry(localReader)
	if err != nil {
		return err
	}
	if local.Path != dirEntry.Path {
		return fmt.Errorf("path mismatch")
	}
	if local.Header.EntryType != dirEntry.EntryType || local.Header.CompressionMethod != dirEntry.CompressionMethod || local.Header.UncompressedSize != dirEntry.UncompressedSize || local.Header.CompressedSize != dirEntry.CompressedSize || local.Header.DataCRC32 != dirEntry.DataCRC32 {
		return fmt.Errorf("metadata mismatch")
	}
	if local.Header.ExtraMetadataLength != 0 {
		return fmt.Errorf("extra metadata is not supported")
	}
	return nil
}
