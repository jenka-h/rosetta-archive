package format

import (
	"encoding/binary"
	"strings"
	"unicode/utf8"

	core "rosetta-archive/internal"
)

const (
	// ArchiveMagic is the only magic value in the format. Entry, directory, and footer
	// structures are located by fixed layout and offsets instead of repeated magic bytes.
	ArchiveMagic = "ROSA"

	FormatMajor uint16 = 1
	FormatMinor uint16 = 0

	ArchiveHeaderSize   uint32 = 16
	EntryHeaderSize     uint32 = 42
	DirectoryHeaderSize uint32 = 12
	DirectoryEntrySize  uint32 = 46
	FooterSize          uint32 = 28

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

func validateEntryType(t EntryType) error {
	switch t {
	case EntryTypeFile, EntryTypeDirectory:
		return nil
	default:
		return core.ErrInvalidEntryType
	}
}

func validateCompressionMethod(m CompressionMethod) error {
	switch m {
	case CompressionStore, CompressionROSA1:
		return nil
	default:
		return core.ErrInvalidCompressionMethod
	}
}

func validatePathBytes(path []byte) error {
	if len(path) == 0 {
		return core.ErrInvalidPath
	}
	if uint32(len(path)) > MaxPathLength {
		return core.ErrPathTooLong
	}
	if !utf8.Valid(path) {
		return core.ErrInvalidPath
	}
	p := string(path)
	if strings.ContainsRune(p, '\x00') || strings.HasPrefix(p, "/") || strings.Contains(p, "\\") {
		return core.ErrInvalidPath
	}
	if len(p) >= 2 && p[1] == ':' && ((p[0] >= 'A' && p[0] <= 'Z') || (p[0] >= 'a' && p[0] <= 'z')) {
		return core.ErrInvalidPath
	}
	parts := strings.Split(p, "/")
	for i, part := range parts {
		if part == "" {
			if i == len(parts)-1 {
				continue
			}
			return core.ErrInvalidPath
		}
		if part == "." || part == ".." {
			return core.ErrInvalidPath
		}
	}
	return nil
}

func putMagic(dst []byte, magic string) {
	copy(dst, []byte(magic))
}

func hasMagic(src []byte, magic string) bool {
	return string(src[:4]) == magic
}
