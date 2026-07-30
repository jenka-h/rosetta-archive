package format

import (
	"fmt"
	"io"
)

// ArchiveHeader is the fixed 16-byte ROSA archive header.
type ArchiveHeader struct {
	Version     uint16
	GlobalFlags uint16
	HeaderSize  uint32
	Reserved    uint32
}

// Header is kept as a compatibility alias for older scaffold call sites.
type Header = ArchiveHeader

func NewArchiveHeader(globalFlags uint16) ArchiveHeader {
	return ArchiveHeader{
		Version:     FormatVersion,
		GlobalFlags: globalFlags,
		HeaderSize:  ArchiveHeaderSize,
	}
}

func EncodeHeader(w io.Writer, h ArchiveHeader) error {
	if w == nil {
		return ErrNilWriter
	}
	if err := ValidateHeader(h); err != nil {
		return fmt.Errorf("validate archive header: %w", err)
	}
	buf := make([]byte, ArchiveHeaderSize)
	putMagic(buf[0:4], ArchiveMagic)
	ByteOrder.PutUint16(buf[4:6], h.Version)
	ByteOrder.PutUint16(buf[6:8], h.GlobalFlags)
	ByteOrder.PutUint32(buf[8:12], h.HeaderSize)
	ByteOrder.PutUint32(buf[12:16], h.Reserved)
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
	h := ArchiveHeader{
		Version:     ByteOrder.Uint16(buf[4:6]),
		GlobalFlags: ByteOrder.Uint16(buf[6:8]),
		HeaderSize:  ByteOrder.Uint32(buf[8:12]),
		Reserved:    ByteOrder.Uint32(buf[12:16]),
	}
	if err := ValidateHeader(h); err != nil {
		return ArchiveHeader{}, fmt.Errorf("validate archive header: %w", err)
	}
	return h, nil
}

func ValidateHeader(h ArchiveHeader) error {
	if h.Version != FormatVersion {
		return ErrUnsupportedVersion
	}
	if !knownUint16Flags(h.GlobalFlags, GlobalFlagsKnownMask) {
		return ErrUnknownFlags
	}
	if h.HeaderSize != ArchiveHeaderSize {
		return ErrInvalidSize
	}
	if h.Reserved != 0 {
		return ErrReservedNonZero
	}
	return nil
}
