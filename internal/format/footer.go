package format

import (
	"fmt"
	"io"

	"rosetta-archive/internal"
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

// ArchiveFooter is the fixed final archive structure. It has no magic and no CRC;
// integrity is covered by HeaderCRC32, DataCRC32, and DirectoryCRC32.
type ArchiveFooter struct {
	MajorVersion           uint16
	MinorVersion           uint16
	CentralDirectoryOffset uint64
	CentralDirectorySize   uint64
	EntryCount             uint32
	FooterSize             uint32
}

type Footer = ArchiveFooter

func EncodeDirectoryHeader(w io.Writer, h DirectoryHeader) error {
	if w == nil {
		return ErrNilWriter
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
		return DirectoryHeader{}, ErrNilReader
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
	ByteOrder.PutUint32(buf[2:6], e.PathLength)
	ByteOrder.PutUint64(buf[6:14], e.FileRecordOffset)
	ByteOrder.PutUint64(buf[14:22], e.DataOffset)
	ByteOrder.PutUint64(buf[22:30], e.UncompressedSize)
	ByteOrder.PutUint64(buf[30:38], e.CompressedSize)
	ByteOrder.PutUint32(buf[38:42], e.DataCRC32)
	ByteOrder.PutUint32(buf[42:46], e.Reserved)
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
		PathLength:        ByteOrder.Uint32(buf[2:6]),
		FileRecordOffset:  ByteOrder.Uint64(buf[6:14]),
		DataOffset:        ByteOrder.Uint64(buf[14:22]),
		UncompressedSize:  ByteOrder.Uint64(buf[22:30]),
		CompressedSize:    ByteOrder.Uint64(buf[30:38]),
		DataCRC32:         ByteOrder.Uint32(buf[38:42]),
		Reserved:          ByteOrder.Uint32(buf[42:46]),
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
	if f.MajorVersion == 0 {
		f.MajorVersion = FormatMajor
	}
	f.FooterSize = FooterSize
	if err := ValidateFooter(f); err != nil {
		return fmt.Errorf("validate archive footer: %w", err)
	}
	if _, err := w.Write(marshalFooter(f)); err != nil {
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
	f := unmarshalFooter(buf)
	if err := ValidateFooter(f); err != nil {
		return ArchiveFooter{}, fmt.Errorf("validate archive footer: %w", err)
	}
	return f, nil
}

func ValidateDirectoryHeader(h DirectoryHeader, verifyCRCField bool) error {
	if h.Reserved != 0 {
		return ErrReservedNonZero
	}
	if !verifyCRCField && h.DirectoryCRC32 != 0 {
		return ErrInvalidDirectoryCRC32
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
	return nil
}

func ValidateFooter(f ArchiveFooter) error {
	if _, err := CheckCompatibility(f.MajorVersion, f.MinorVersion); err != nil {
		return err
	}
	if f.FooterSize != FooterSize {
		return ErrInvalidSize
	}
	return nil
}

func DirectoryCRC32(directory []byte) uint32 {
	if len(directory) >= 8 {
		copyForCRC := make([]byte, len(directory))
		copy(copyForCRC, directory)
		ByteOrder.PutUint32(copyForCRC[4:8], 0)
		return internal.CRC32(copyForCRC)
	}
	return internal.CRC32(directory)
}

func marshalDirectoryHeader(h DirectoryHeader) []byte {
	buf := make([]byte, DirectoryHeaderSize)
	ByteOrder.PutUint32(buf[0:4], h.EntryCount)
	ByteOrder.PutUint32(buf[4:8], h.DirectoryCRC32)
	ByteOrder.PutUint32(buf[8:12], h.Reserved)
	return buf
}

func unmarshalDirectoryHeader(buf []byte) DirectoryHeader {
	return DirectoryHeader{
		EntryCount:     ByteOrder.Uint32(buf[0:4]),
		DirectoryCRC32: ByteOrder.Uint32(buf[4:8]),
		Reserved:       ByteOrder.Uint32(buf[8:12]),
	}
}

func marshalFooter(f ArchiveFooter) []byte {
	buf := make([]byte, FooterSize)
	ByteOrder.PutUint16(buf[0:2], f.MajorVersion)
	ByteOrder.PutUint16(buf[2:4], f.MinorVersion)
	ByteOrder.PutUint64(buf[4:12], f.CentralDirectoryOffset)
	ByteOrder.PutUint64(buf[12:20], f.CentralDirectorySize)
	ByteOrder.PutUint32(buf[20:24], f.EntryCount)
	ByteOrder.PutUint32(buf[24:28], f.FooterSize)
	return buf
}

func unmarshalFooter(buf []byte) ArchiveFooter {
	return ArchiveFooter{
		MajorVersion:           ByteOrder.Uint16(buf[0:2]),
		MinorVersion:           ByteOrder.Uint16(buf[2:4]),
		CentralDirectoryOffset: ByteOrder.Uint64(buf[4:12]),
		CentralDirectorySize:   ByteOrder.Uint64(buf[12:20]),
		EntryCount:             ByteOrder.Uint32(buf[20:24]),
		FooterSize:             ByteOrder.Uint32(buf[24:28]),
	}
}
