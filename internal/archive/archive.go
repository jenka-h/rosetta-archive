package archive

import (
	"io"

	"rosetta-archive/internal/format"
)

type Entry struct {
	Path       string
	Type       format.EntryType
	Size       uint64
	CRC32      uint32
	Offset     uint64
	DataOffset uint64
}

type InfoSummary struct {
	FormatName   string
	MajorVersion uint8
	MinorVersion uint8
	ArchiveSize  uint64

	EntryCount     uint32
	FileCount      uint32
	DirectoryCount uint32

	StoredSize       uint64
	UncompressedSize uint64
}

func closeIfNeeded(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}

func entryFromDirectory(e format.DirectoryEntry) Entry {
	return Entry{
		Path:       e.Path,
		Type:       e.EntryType,
		Size:       e.UncompressedSize,
		CRC32:      e.DataCRC32,
		Offset:     e.FileRecordOffset,
		DataOffset: e.DataOffset,
	}
}
