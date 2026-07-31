package format

import (
	"fmt"
	"io"

	"rosetta-archive/internal"
)

// ArchiveHeader is the fixed archive header.
//
// Layout, little-endian:
//   - 0..3:   magic "ROSA"
//   - 4..5:   major version
//   - 6..7:   minor version
//   - 8..11:  header size
//   - 12..15: header CRC-32, calculated with this field set to zero
//
// This is the only structure that contains a magic value.
type ArchiveHeader struct {
	MajorVersion uint16
	MinorVersion uint16
	HeaderSize   uint32
	HeaderCRC32  uint32
}

type Header = ArchiveHeader

func NewArchiveHeader() ArchiveHeader {
	return ArchiveHeader{
		MajorVersion: FormatMajor,
		MinorVersion: FormatMinor,
		HeaderSize:   ArchiveHeaderSize,
	}
}

func EncodeHeader(w io.Writer, h ArchiveHeader) error {
	if w == nil {
		return ErrNilWriter
	}
	if h.MajorVersion == 0 {
		h.MajorVersion = FormatMajor
	}
	h.HeaderSize = ArchiveHeaderSize
	h.HeaderCRC32 = 0
	if err := ValidateHeader(h, false); err != nil {
		return fmt.Errorf("validate archive header: %w", err)
	}
	buf := marshalHeader(h)
	ByteOrder.PutUint32(buf[12:16], internal.CRC32(buf))
	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("write archive header: %w", err)
	}
	return nil
}

func DecodeHeader(r io.Reader) (ArchiveHeader, error) {
	if r == nil {
		return ArchiveHeader{}, ErrNilReader
	}
	buf := make([]byte, ArchiveHeaderSize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return ArchiveHeader{}, fmt.Errorf("read archive header: %w", err)
	}
	if !hasMagic(buf, ArchiveMagic) {
		return ArchiveHeader{}, ErrInvalidMagic
	}
	h := unmarshalHeader(buf)
	storedCRC := h.HeaderCRC32
	ByteOrder.PutUint32(buf[12:16], 0)
	if internal.CRC32(buf) != storedCRC {
		return ArchiveHeader{}, ErrInvalidHeaderCRC32
	}
	if err := ValidateHeader(h, true); err != nil {
		return ArchiveHeader{}, fmt.Errorf("validate archive header: %w", err)
	}
	return h, nil
}

func ValidateHeader(h ArchiveHeader, verifyCRCField bool) error {
	if _, err := CheckCompatibility(h.MajorVersion, h.MinorVersion); err != nil {
		return err
	}
	if h.HeaderSize != ArchiveHeaderSize {
		return ErrInvalidSize
	}
	if !verifyCRCField && h.HeaderCRC32 != 0 {
		return ErrInvalidHeaderCRC32
	}
	return nil
}

func marshalHeader(h ArchiveHeader) []byte {
	buf := make([]byte, ArchiveHeaderSize)
	putMagic(buf[0:4], ArchiveMagic)
	ByteOrder.PutUint16(buf[4:6], h.MajorVersion)
	ByteOrder.PutUint16(buf[6:8], h.MinorVersion)
	ByteOrder.PutUint32(buf[8:12], h.HeaderSize)
	ByteOrder.PutUint32(buf[12:16], h.HeaderCRC32)
	return buf
}

func unmarshalHeader(buf []byte) ArchiveHeader {
	return ArchiveHeader{
		MajorVersion: ByteOrder.Uint16(buf[4:6]),
		MinorVersion: ByteOrder.Uint16(buf[6:8]),
		HeaderSize:   ByteOrder.Uint32(buf[8:12]),
		HeaderCRC32:  ByteOrder.Uint32(buf[12:16]),
	}
}
