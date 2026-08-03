package format

import (
	"fmt"
	"io"

	core "rosetta-archive/internal"
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

func NewArchiveHeader() ArchiveHeader {
	return ArchiveHeader{
		MajorVersion: FormatMajor,
		MinorVersion: FormatMinor,
		HeaderSize:   ArchiveHeaderSize,
	}
}

func EncodeHeader(w io.Writer, h ArchiveHeader) error {
	if w == nil {
		return core.ErrNilWriter
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
	ByteOrder.PutUint32(buf[12:16], core.CRC32(buf))
	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("write archive header: %w", err)
	}
	return nil
}

func DecodeHeader(r io.Reader) (ArchiveHeader, error) {
	if r == nil {
		return ArchiveHeader{}, core.ErrNilReader
	}
	buf := make([]byte, ArchiveHeaderSize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return ArchiveHeader{}, fmt.Errorf("read archive header: %w", err)
	}
	if !hasMagic(buf, ArchiveMagic) {
		return ArchiveHeader{}, core.ErrInvalidMagic
	}
	h := unmarshalHeader(buf)
	storedCRC := h.HeaderCRC32
	ByteOrder.PutUint32(buf[12:16], 0)
	if core.CRC32(buf) != storedCRC {
		return ArchiveHeader{}, core.ErrInvalidHeaderCRC32
	}
	if err := ValidateHeader(h, true); err != nil {
		return ArchiveHeader{}, fmt.Errorf("validate archive header: %w", err)
	}
	return h, nil
}

func ValidateHeader(h ArchiveHeader, verifyCRCField bool) error {
	if err := CheckCompatibility(h.MajorVersion, h.MinorVersion); err != nil {
		return err
	}
	if h.HeaderSize != ArchiveHeaderSize {
		return core.ErrInvalidSize
	}
	if !verifyCRCField && h.HeaderCRC32 != 0 {
		return core.ErrInvalidHeaderCRC32
	}
	return nil
}
