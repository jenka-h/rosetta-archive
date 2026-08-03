package format

import "encoding/binary"

const (
	ArchiveMagic = "ROSA"

	FormatMajor uint16 = 1
	FormatMinor uint16 = 0

	ArchiveHeaderSize   uint32 = 16
	EntryHeaderSize     uint32 = 42
	DirectoryHeaderSize uint32 = 12
	DirectoryEntrySize  uint32 = 46
	FooterSize          uint32 = 24

	MaxPathLength          uint32 = 1 << 20
	MaxExtraMetadataLength uint32 = 16 << 20
)

// ByteOrder is the byte order used for every multibyte field in ROSA v1.
var ByteOrder = binary.LittleEndian

type EntryType uint8

const (
	EntryTypeFile EntryType = iota
	EntryTypeDirectory
)

type CompressionMethod uint8

const (
	CompressionStore CompressionMethod = iota
	CompressionROSA1
)
