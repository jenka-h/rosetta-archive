package format

import (
	"fmt"
	"io"

	core "rosetta-archive/internal"
)

// DirectoryHeader begins the central directory.
// DirectoryCRC32 is calculated over the complete central directory bytes with this
// field encoded as zero.
type DirectoryHeader struct {
	EntryCount     uint32
	DirectoryCRC32 uint32
	Reserved       uint32
}

// DirectoryEntry is one fixed central-directory record plus its UTF-8 path.
type DirectoryEntry struct {
	EntryType         EntryType
	CompressionMethod CompressionMethod
	PathLength        uint32
	FileRecordOffset  uint64
	DataOffset        uint64
	UncompressedSize  uint64
	CompressedSize    uint64
	DataCRC32         uint32
	Reserved          uint32
	Path              string
}

func EncodeDirectoryHeader(w io.Writer, h DirectoryHeader) error {
	if w == nil {
		return core.ErrNilWriter
	}
	if err := ValidateDirectoryHeader(h, false); err != nil {
		return fmt.Errorf("validate directory header: %w", err)
	}
	buf := marshalDirectoryHeader(h)
	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("write directory header: %w", err)
	}
	return nil
}

func DecodeDirectoryHeader(r io.Reader) (DirectoryHeader, error) {
	if r == nil {
		return DirectoryHeader{}, core.ErrNilReader
	}
	buf := make([]byte, DirectoryHeaderSize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return DirectoryHeader{}, fmt.Errorf("read directory header: %w", err)
	}
	h := unmarshalDirectoryHeader(buf)
	if err := ValidateDirectoryHeader(h, true); err != nil {
		return DirectoryHeader{}, fmt.Errorf("validate directory header: %w", err)
	}
	return h, nil
}

func EncodeDirectoryEntry(w io.Writer, e DirectoryEntry) error {
	if w == nil {
		return core.ErrNilWriter
	}
	pathBytes := []byte(e.Path)
	if err := validatePathBytes(pathBytes); err != nil {
		return fmt.Errorf("validate directory entry path: %w", err)
	}
	e.PathLength = uint32(len(pathBytes))
	if err := ValidateDirectoryEntry(e); err != nil {
		return fmt.Errorf("validate directory entry: %w", err)
	}
	if _, err := w.Write(marshalDirectoryEntry(e)); err != nil {
		return fmt.Errorf("write directory entry: %w", err)
	}
	if _, err := w.Write(pathBytes); err != nil {
		return fmt.Errorf("write directory entry path: %w", err)
	}
	return nil
}

func DecodeDirectoryEntry(r io.Reader) (DirectoryEntry, error) {
	if r == nil {
		return DirectoryEntry{}, core.ErrNilReader
	}
	buf := make([]byte, DirectoryEntrySize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return DirectoryEntry{}, fmt.Errorf("read directory entry: %w", err)
	}
	e := unmarshalDirectoryEntry(buf)
	if e.PathLength == 0 || e.PathLength > MaxPathLength {
		return DirectoryEntry{}, core.ErrPathTooLong
	}
	pathBytes := make([]byte, e.PathLength)
	if _, err := io.ReadFull(r, pathBytes); err != nil {
		return DirectoryEntry{}, fmt.Errorf("read directory entry path: %w", err)
	}
	if err := validatePathBytes(pathBytes); err != nil {
		return DirectoryEntry{}, fmt.Errorf("validate directory entry path: %w", err)
	}
	e.Path = string(pathBytes)
	if err := ValidateDirectoryEntry(e); err != nil {
		return DirectoryEntry{}, fmt.Errorf("validate directory entry: %w", err)
	}
	return e, nil
}

func ValidateDirectoryHeader(h DirectoryHeader, verifyCRCField bool) error {
	if h.Reserved != 0 {
		return core.ErrReservedNonZero
	}
	if !verifyCRCField && h.DirectoryCRC32 != 0 {
		return core.ErrInvalidDirectoryCRC32
	}
	return nil
}

func ValidateDirectoryEntry(e DirectoryEntry) error {
	if err := validateEntryType(e.EntryType); err != nil {
		return err
	}
	if err := validateCompressionMethod(e.CompressionMethod); err != nil {
		return err
	}
	if e.PathLength == 0 || e.PathLength > MaxPathLength || int(e.PathLength) != len([]byte(e.Path)) {
		return core.ErrPathTooLong
	}
	if err := validatePathBytes([]byte(e.Path)); err != nil {
		return err
	}
	if e.Reserved != 0 {
		return core.ErrReservedNonZero
	}
	if e.EntryType == EntryTypeDirectory {
		if e.CompressionMethod != CompressionStore || e.UncompressedSize != 0 || e.CompressedSize != 0 || e.DataCRC32 != 0 {
			return core.ErrInvalidEntrySizes
		}
	}
	if e.CompressionMethod == CompressionStore && e.EntryType == EntryTypeFile && e.CompressedSize != e.UncompressedSize {
		return core.ErrInvalidEntrySizes
	}
	return nil
}

func DirectoryCRC32(directory []byte) uint32 {
	if len(directory) < int(DirectoryHeaderSize) {
		return core.CRC32(directory)
	}
	copyForCRC := make([]byte, len(directory))
	copy(copyForCRC, directory)
	ByteOrder.PutUint32(copyForCRC[4:8], 0)
	return core.CRC32(copyForCRC)
}
