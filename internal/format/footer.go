package format

import (
	"fmt"
	"io"

	"rosetta-archive/internal"
)

// DirectoryHeader is the fixed 16-byte header that begins the central directory.
type DirectoryHeader struct {
	EntryCount     uint32
	DirectoryFlags uint32
	Reserved       uint32
}

// DirectoryEntry is one fixed 48-byte central-directory record plus its UTF-8 path.
type DirectoryEntry struct {
	EntryType         EntryType
	CompressionMethod CompressionMethod
	EntryFlags        uint16
	PathLength        uint32
	FileRecordOffset  uint64
	DataOffset        uint64
	UncompressedSize  uint64
	CompressedSize    uint64
	DataCRC32         uint32
	Reserved          uint32
	Path              string
}

// ArchiveFooter is the fixed 40-byte final archive structure.
type ArchiveFooter struct {
	Version                uint16
	FooterFlags            uint16
	CentralDirectoryOffset uint64
	CentralDirectorySize   uint64
	EntryCount             uint32
	CentralDirectoryCRC32  uint32
	FooterCRC32            uint32
	FooterSize             uint32
}

// Footer is kept as a compatibility alias for older scaffold call sites.
type Footer = ArchiveFooter

func EncodeDirectoryHeader(w io.Writer, h DirectoryHeader) error {
	if w == nil {
		return ErrNilWriter
	}
	if err := ValidateDirectoryHeader(h); err != nil {
		return fmt.Errorf("validate directory header: %w", err)
	}
	buf := make([]byte, DirectoryHeaderSize)
	putMagic(buf[0:4], DirectoryMagic)
	ByteOrder.PutUint32(buf[4:8], h.EntryCount)
	ByteOrder.PutUint32(buf[8:12], h.DirectoryFlags)
	ByteOrder.PutUint32(buf[12:16], h.Reserved)
	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("write directory header: %w", err)
	}
	return nil
}

func DecodeDirectoryHeader(r io.Reader) (DirectoryHeader, error) {
	if r == nil {
		return DirectoryHeader{}, ErrNilReader
	}
	buf := make([]byte, DirectoryHeaderSize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return DirectoryHeader{}, fmt.Errorf("read directory header: %w", err)
	}
	if !hasMagic(buf, DirectoryMagic) {
		return DirectoryHeader{}, ErrInvalidMagic
	}
	h := DirectoryHeader{
		EntryCount:     ByteOrder.Uint32(buf[4:8]),
		DirectoryFlags: ByteOrder.Uint32(buf[8:12]),
		Reserved:       ByteOrder.Uint32(buf[12:16]),
	}
	if err := ValidateDirectoryHeader(h); err != nil {
		return DirectoryHeader{}, fmt.Errorf("validate directory header: %w", err)
	}
	return h, nil
}

func EncodeDirectoryEntry(w io.Writer, e DirectoryEntry) error {
	if w == nil {
		return ErrNilWriter
	}
	pathBytes := []byte(e.Path)
	if err := validatePathBytes(pathBytes); err != nil {
		return fmt.Errorf("validate directory entry path: %w", err)
	}
	e.PathLength = uint32(len(pathBytes))
	if err := ValidateDirectoryEntry(e); err != nil {
		return fmt.Errorf("validate directory entry: %w", err)
	}
	buf := make([]byte, DirectoryEntrySize)
	buf[0] = byte(e.EntryType)
	buf[1] = byte(e.CompressionMethod)
	ByteOrder.PutUint16(buf[2:4], e.EntryFlags)
	ByteOrder.PutUint32(buf[4:8], e.PathLength)
	ByteOrder.PutUint64(buf[8:16], e.FileRecordOffset)
	ByteOrder.PutUint64(buf[16:24], e.DataOffset)
	ByteOrder.PutUint64(buf[24:32], e.UncompressedSize)
	ByteOrder.PutUint64(buf[32:40], e.CompressedSize)
	ByteOrder.PutUint32(buf[40:44], e.DataCRC32)
	ByteOrder.PutUint32(buf[44:48], e.Reserved)
	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("write directory entry: %w", err)
	}
	if _, err := w.Write(pathBytes); err != nil {
		return fmt.Errorf("write directory entry path: %w", err)
	}
	return nil
}

func DecodeDirectoryEntry(r io.Reader) (DirectoryEntry, error) {
	if r == nil {
		return DirectoryEntry{}, ErrNilReader
	}
	buf := make([]byte, DirectoryEntrySize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return DirectoryEntry{}, fmt.Errorf("read directory entry: %w", err)
	}
	e := DirectoryEntry{
		EntryType:         EntryType(buf[0]),
		CompressionMethod: CompressionMethod(buf[1]),
		EntryFlags:        ByteOrder.Uint16(buf[2:4]),
		PathLength:        ByteOrder.Uint32(buf[4:8]),
		FileRecordOffset:  ByteOrder.Uint64(buf[8:16]),
		DataOffset:        ByteOrder.Uint64(buf[16:24]),
		UncompressedSize:  ByteOrder.Uint64(buf[24:32]),
		CompressedSize:    ByteOrder.Uint64(buf[32:40]),
		DataCRC32:         ByteOrder.Uint32(buf[40:44]),
		Reserved:          ByteOrder.Uint32(buf[44:48]),
	}
	if e.PathLength == 0 || e.PathLength > MaxPathLength {
		return DirectoryEntry{}, ErrPathTooLong
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

func EncodeFooter(w io.Writer, f ArchiveFooter) error {
	if w == nil {
		return ErrNilWriter
	}
	f.FooterSize = FooterSize
	if f.Version == 0 {
		f.Version = FormatVersion
	}
	if f.FooterFlags&FooterFlagCentralDirectoryCRC == 0 {
		f.CentralDirectoryCRC32 = 0
	}
	f.FooterCRC32 = 0
	if err := ValidateFooter(f, false); err != nil {
		return fmt.Errorf("validate archive footer: %w", err)
	}
	buf := marshalFooter(f)
	crc := internal.CRC32(buf)
	ByteOrder.PutUint32(buf[32:36], crc)
	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("write archive footer: %w", err)
	}
	return nil
}

func DecodeFooter(r io.Reader) (ArchiveFooter, error) {
	if r == nil {
		return ArchiveFooter{}, ErrNilReader
	}
	buf := make([]byte, FooterSize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return ArchiveFooter{}, fmt.Errorf("read archive footer: %w", err)
	}
	if !hasMagic(buf, FooterMagic) {
		return ArchiveFooter{}, ErrInvalidMagic
	}
	f := unmarshalFooter(buf)
	storedCRC := f.FooterCRC32
	ByteOrder.PutUint32(buf[32:36], 0)
	if internal.CRC32(buf) != storedCRC {
		return ArchiveFooter{}, ErrInvalidFooterCRC32
	}
	if err := ValidateFooter(f, true); err != nil {
		return ArchiveFooter{}, fmt.Errorf("validate archive footer: %w", err)
	}
	return f, nil
}

func ValidateDirectoryHeader(h DirectoryHeader) error {
	if !knownUint32Flags(h.DirectoryFlags, DirectoryFlagsKnownMask) {
		return ErrUnknownFlags
	}
	if h.Reserved != 0 {
		return ErrReservedNonZero
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
	if !knownUint16Flags(e.EntryFlags, EntryFlagsKnownMask) {
		return ErrUnknownFlags
	}
	if e.PathLength == 0 || e.PathLength > MaxPathLength || int(e.PathLength) != len([]byte(e.Path)) {
		return ErrPathTooLong
	}
	if err := validatePathBytes([]byte(e.Path)); err != nil {
		return err
	}
	if e.Reserved != 0 {
		return ErrReservedNonZero
	}
	if e.EntryType == EntryTypeDirectory {
		if e.CompressionMethod != CompressionStore || e.UncompressedSize != 0 || e.CompressedSize != 0 || e.DataCRC32 != 0 {
			return ErrInvalidEntrySizes
		}
	}
	if e.CompressionMethod == CompressionStore && e.EntryType == EntryTypeFile && e.CompressedSize != e.UncompressedSize {
		return ErrInvalidEntrySizes
	}
	if e.EntryFlags&EntryFlagCRC32Present == 0 && e.DataCRC32 != 0 {
		return ErrInvalidEntryCRC
	}
	return nil
}

func ValidateFooter(f ArchiveFooter, verifyCRCField bool) error {
	if f.Version != FormatVersion {
		return ErrUnsupportedVersion
	}
	if !knownUint16Flags(f.FooterFlags, FooterFlagsKnownMask) {
		return ErrUnknownFlags
	}
	if f.FooterFlags&FooterFlagCentralDirectoryCRC == 0 && f.CentralDirectoryCRC32 != 0 {
		return ErrInvalidEntryCRC
	}
	if f.FooterSize != FooterSize {
		return ErrInvalidSize
	}
	if !verifyCRCField && f.FooterCRC32 != 0 {
		return ErrInvalidFooterCRC32
	}
	return nil
}

func marshalFooter(f ArchiveFooter) []byte {
	buf := make([]byte, FooterSize)
	putMagic(buf[0:4], FooterMagic)
	ByteOrder.PutUint16(buf[4:6], f.Version)
	ByteOrder.PutUint16(buf[6:8], f.FooterFlags)
	ByteOrder.PutUint64(buf[8:16], f.CentralDirectoryOffset)
	ByteOrder.PutUint64(buf[16:24], f.CentralDirectorySize)
	ByteOrder.PutUint32(buf[24:28], f.EntryCount)
	ByteOrder.PutUint32(buf[28:32], f.CentralDirectoryCRC32)
	ByteOrder.PutUint32(buf[32:36], f.FooterCRC32)
	ByteOrder.PutUint32(buf[36:40], f.FooterSize)
	return buf
}

func unmarshalFooter(buf []byte) ArchiveFooter {
	return ArchiveFooter{
		Version:                ByteOrder.Uint16(buf[4:6]),
		FooterFlags:            ByteOrder.Uint16(buf[6:8]),
		CentralDirectoryOffset: ByteOrder.Uint64(buf[8:16]),
		CentralDirectorySize:   ByteOrder.Uint64(buf[16:24]),
		EntryCount:             ByteOrder.Uint32(buf[24:28]),
		CentralDirectoryCRC32:  ByteOrder.Uint32(buf[28:32]),
		FooterCRC32:            ByteOrder.Uint32(buf[32:36]),
		FooterSize:             ByteOrder.Uint32(buf[36:40]),
	}
}
