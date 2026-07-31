package format

import (
	"fmt"
	"io"
)

// EntryHeader is the fixed header that begins every file record.
// It contains no magic and no flags; DataCRC32 is always the CRC-32 field for file data.
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
		return ErrNilWriter
	}
	if err := ValidateEntryHeader(h); err != nil {
		return fmt.Errorf("validate entry header: %w", err)
	}
	buf := make([]byte, EntryHeaderSize)
	buf[0] = byte(h.EntryType)
	buf[1] = byte(h.CompressionMethod)
	ByteOrder.PutUint32(buf[2:6], h.PathLength)
	ByteOrder.PutUint32(buf[6:10], h.ExtraMetadataLength)
	ByteOrder.PutUint64(buf[10:18], h.UncompressedSize)
	ByteOrder.PutUint64(buf[18:26], h.CompressedSize)
	ByteOrder.PutUint64(buf[26:34], uint64(h.ModificationTime))
	ByteOrder.PutUint32(buf[34:38], h.DataCRC32)
	ByteOrder.PutUint32(buf[38:42], h.Reserved)
	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("write entry header: %w", err)
	}
	return nil
}

func DecodeEntryHeader(r io.Reader) (EntryHeader, error) {
	if r == nil {
		return EntryHeader{}, ErrNilReader
	}
	buf := make([]byte, EntryHeaderSize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return EntryHeader{}, fmt.Errorf("read entry header: %w", err)
	}
	h := EntryHeader{
		EntryType:           EntryType(buf[0]),
		CompressionMethod:   CompressionMethod(buf[1]),
		PathLength:          ByteOrder.Uint32(buf[2:6]),
		ExtraMetadataLength: ByteOrder.Uint32(buf[6:10]),
		UncompressedSize:    ByteOrder.Uint64(buf[10:18]),
		CompressedSize:      ByteOrder.Uint64(buf[18:26]),
		ModificationTime:    int64(ByteOrder.Uint64(buf[26:34])),
		DataCRC32:           ByteOrder.Uint32(buf[34:38]),
		Reserved:            ByteOrder.Uint32(buf[38:42]),
	}
	if err := ValidateEntryHeader(h); err != nil {
		return EntryHeader{}, fmt.Errorf("validate entry header: %w", err)
	}
	return h, nil
}

func EncodeEntry(w io.Writer, e Entry) error {
	if w == nil {
		return ErrNilWriter
	}
	pathBytes := []byte(e.Path)
	if err := validatePathBytes(pathBytes); err != nil {
		return fmt.Errorf("validate entry path: %w", err)
	}
	if len(e.ExtraMetadata) > int(MaxExtraMetadataLength) {
		return ErrExtraMetadataTooLong
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
		return Entry{}, ErrNilReader
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
		return ErrPathTooLong
	}
	if h.ExtraMetadataLength > MaxExtraMetadataLength {
		return ErrExtraMetadataTooLong
	}
	if h.Reserved != 0 {
		return ErrReservedNonZero
	}
	if h.EntryType == EntryTypeDirectory {
		if h.CompressionMethod != CompressionStore || h.UncompressedSize != 0 || h.CompressedSize != 0 || h.DataCRC32 != 0 {
			return ErrInvalidEntrySizes
		}
	}
	if h.CompressionMethod == CompressionStore && h.EntryType == EntryTypeFile && h.CompressedSize != h.UncompressedSize {
		return ErrInvalidEntrySizes
	}
	return nil
}
