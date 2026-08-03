package format

import (
	"fmt"
	"io"

	core "rosetta-archive/internal"
)

type EntryHeader struct {
	EntryType           EntryType
	CompressionMethod   CompressionMethod
	PathLength          uint32
	ExtraMetadataLength uint32
	UncompressedSize    uint64
	CompressedSize      uint64
	ModificationTime    int64
	DataCRC32           uint32
	Reserved            uint32
}

// Entry is a logical file record header plus its path and extra metadata bytes.
type Entry struct {
	Header        EntryHeader
	Path          string
	ExtraMetadata []byte
}

func EncodeEntryHeader(w io.Writer, h EntryHeader) error {
	if w == nil {
		return core.ErrNilWriter
	}
	if err := ValidateEntryHeader(h); err != nil {
		return fmt.Errorf("validate entry header: %w", err)
	}
	buf := marshalEntryHeader(h)
	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("write entry header: %w", err)
	}
	return nil
}

func DecodeEntryHeader(r io.Reader) (EntryHeader, error) {
	if r == nil {
		return EntryHeader{}, core.ErrNilReader
	}
	buf := make([]byte, EntryHeaderSize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return EntryHeader{}, fmt.Errorf("read entry header: %w", err)
	}
	h := unmarshalEntryHeader(buf)
	if err := ValidateEntryHeader(h); err != nil {
		return EntryHeader{}, fmt.Errorf("validate entry header: %w", err)
	}
	return h, nil
}

func EncodeEntry(w io.Writer, e Entry) error {
	if w == nil {
		return core.ErrNilWriter
	}
	pathBytes := []byte(e.Path)
	if err := validatePathBytes(pathBytes); err != nil {
		return fmt.Errorf("validate entry path: %w", err)
	}
	if len(e.ExtraMetadata) > int(MaxExtraMetadataLength) {
		return core.ErrExtraMetadataTooLong
	}
	e.Header.PathLength = uint32(len(pathBytes))
	e.Header.ExtraMetadataLength = uint32(len(e.ExtraMetadata))
	if err := EncodeEntryHeader(w, e.Header); err != nil {
		return err
	}
	if _, err := w.Write(pathBytes); err != nil {
		return fmt.Errorf("write entry path: %w", err)
	}
	if len(e.ExtraMetadata) > 0 {
		if _, err := w.Write(e.ExtraMetadata); err != nil {
			return fmt.Errorf("write entry extra metadata: %w", err)
		}
	}
	return nil
}

func DecodeEntry(r io.Reader) (Entry, error) {
	if r == nil {
		return Entry{}, core.ErrNilReader
	}
	h, err := DecodeEntryHeader(r)
	if err != nil {
		return Entry{}, err
	}
	pathBytes := make([]byte, h.PathLength)
	if _, err := io.ReadFull(r, pathBytes); err != nil {
		return Entry{}, fmt.Errorf("read entry path: %w", err)
	}
	if err := validatePathBytes(pathBytes); err != nil {
		return Entry{}, fmt.Errorf("validate entry path: %w", err)
	}
	extra := make([]byte, h.ExtraMetadataLength)
	if _, err := io.ReadFull(r, extra); err != nil {
		return Entry{}, fmt.Errorf("read entry extra metadata: %w", err)
	}
	return Entry{Header: h, Path: string(pathBytes), ExtraMetadata: extra}, nil
}

func ValidateEntryHeader(h EntryHeader) error {
	if err := validateEntryType(h.EntryType); err != nil {
		return err
	}
	if err := validateCompressionMethod(h.CompressionMethod); err != nil {
		return err
	}
	if h.PathLength == 0 || h.PathLength > MaxPathLength {
		return core.ErrPathTooLong
	}
	if h.ExtraMetadataLength > MaxExtraMetadataLength {
		return core.ErrExtraMetadataTooLong
	}
	if h.Reserved != 0 {
		return core.ErrReservedNonZero
	}
	if h.EntryType == EntryTypeDirectory {
		if h.CompressionMethod != CompressionStore || h.UncompressedSize != 0 || h.CompressedSize != 0 || h.DataCRC32 != 0 {
			return core.ErrInvalidEntrySizes
		}
	}
	if h.CompressionMethod == CompressionStore && h.EntryType == EntryTypeFile && h.CompressedSize != h.UncompressedSize {
		return core.ErrInvalidEntrySizes
	}
	return nil
}
