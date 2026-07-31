package format

import (
	"strings"
	"unicode/utf8"

	core "rosetta-archive/internal"
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
	if hasDrivePrefix(p) {
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

func hasDrivePrefix(p string) bool {
	return len(p) >= 2 && p[1] == ':' && ((p[0] >= 'A' && p[0] <= 'Z') || (p[0] >= 'a' && p[0] <= 'z'))
}

func putMagic(dst []byte, magic string) {
	copy(dst, []byte(magic))
}

func hasMagic(src []byte, magic string) bool {
	return string(src[:4]) == magic
}
