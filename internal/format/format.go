package format

import (
	"encoding/binary"
	"strings"
	"unicode/utf8"
)

const (
	ArchiveMagic   = "ROSA"
	EntryMagic     = "RENT"
	DirectoryMagic = "RDIR"
	FooterMagic    = "REOF"

	FormatVersion uint16 = 1

	ArchiveHeaderSize   uint32 = 16
	EntryHeaderSize     uint32 = 48
	DirectoryHeaderSize uint32 = 16
	DirectoryEntrySize  uint32 = 48
	FooterSize          uint32 = 40

	MaxPathLength          uint32 = 1 << 20
	MaxExtraMetadataLength uint32 = 16 << 20
)

const (
	GlobalFlagDeterministic       uint16 = 1 << 0
	GlobalFlagExplicitDirectories uint16 = 1 << 1
	GlobalFlagsKnownMask          uint16 = GlobalFlagDeterministic | GlobalFlagExplicitDirectories
	EntryFlagCRC32Present         uint16 = 1 << 0
	EntryFlagModificationTime     uint16 = 1 << 1
	EntryFlagExecutable           uint16 = 1 << 2
	EntryFlagsKnownMask           uint16 = EntryFlagCRC32Present | EntryFlagModificationTime | EntryFlagExecutable
	DirectoryFlagSortedByPath     uint32 = 1 << 0
	DirectoryFlagsKnownMask       uint32 = DirectoryFlagSortedByPath
	FooterFlagCentralDirectoryCRC uint16 = 1 << 0
	FooterFlagsKnownMask          uint16 = FooterFlagCentralDirectoryCRC
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
		return ErrInvalidEntryType
	}
}

func validateCompressionMethod(m CompressionMethod) error {
	switch m {
	case CompressionStore, CompressionROSA1:
		return nil
	default:
		return ErrInvalidCompressionMethod
	}
}

func validatePathBytes(path []byte) error {
	if len(path) == 0 {
		return ErrInvalidPath
	}
	if uint32(len(path)) > MaxPathLength {
		return ErrPathTooLong
	}
	if !utf8.Valid(path) {
		return ErrInvalidPath
	}
	p := string(path)
	if strings.ContainsRune(p, '\x00') || strings.HasPrefix(p, "/") || strings.Contains(p, "\\") {
		return ErrInvalidPath
	}
	if len(p) >= 2 && p[1] == ':' && ((p[0] >= 'A' && p[0] <= 'Z') || (p[0] >= 'a' && p[0] <= 'z')) {
		return ErrInvalidPath
	}
	parts := strings.Split(p, "/")
	for i, part := range parts {
		if part == "" {
			if i == len(parts)-1 {
				continue
			}
			return ErrInvalidPath
		}
		if part == "." || part == ".." {
			return ErrInvalidPath
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

func knownUint16Flags(flags uint16, mask uint16) bool {
	return flags&^mask == 0
}

func knownUint32Flags(flags uint32, mask uint32) bool {
	return flags&^mask == 0
}
