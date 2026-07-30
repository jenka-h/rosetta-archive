package format

import (
	"fmt"
	"io"
)

// EntryHeader is the fixed 48-byte header that begins every file record.
type EntryHeader struct {
	EntryType           EntryType
	CompressionMethod   CompressionMethod
	EntryFlags          uint16
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
	putMagic(buf[0:4], EntryMagic)
	buf[4] = byte(h.EntryType)
	buf[5] = byte(h.CompressionMethod)
	ByteOrder.PutUint16(buf[6:8], h.EntryFlags)
	ByteOrder.PutUint32(buf[8:12], h.PathLength)
	ByteOrder.PutUint32(buf[12:16], h.ExtraMetadataLength)
	ByteOrder.PutUint64(buf[16:24], h.UncompressedSize)
	ByteOrder.PutUint64(buf[24:32], h.CompressedSize)
	ByteOrder.PutUint64(buf[32:40], uint64(h.ModificationTime))
	ByteOrder.PutUint32(buf[40:44], h.DataCRC32)
	ByteOrder.PutUint32(buf[44:48], h.Reserved)
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
	if !hasMagic(buf, EntryMagic) {
		return EntryHeader{}, ErrInvalidMagic
	}
	h := EntryHeader{
		EntryType:           EntryType(buf[4]),
		CompressionMethod:   CompressionMethod(buf[5]),
		EntryFlags:          ByteOrder.Uint16(buf[6:8]),
		PathLength:          ByteOrder.Uint32(buf[8:12]),
		ExtraMetadataLength: ByteOrder.Uint32(buf[12:16]),
		UncompressedSize:    ByteOrder.Uint64(buf[16:24]),
		CompressedSize:      ByteOrder.Uint64(buf[24:32]),
		ModificationTime:    int64(ByteOrder.Uint64(buf[32:40])),
		DataCRC32:           ByteOrder.Uint32(buf[40:44]),
		Reserved:            ByteOrder.Uint32(buf[44:48]),
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
	if !knownUint16Flags(h.EntryFlags, EntryFlagsKnownMask) {
		return ErrUnknownFlags
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
	if h.EntryFlags&EntryFlagCRC32Present == 0 && h.DataCRC32 != 0 {
		return ErrInvalidEntryCRC
	}
	if h.EntryFlags&EntryFlagModificationTime == 0 && h.ModificationTime != 0 {
		return ErrInvalidModificationTime
	}
	return nil
}
