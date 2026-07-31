package format

import (
	"fmt"
	"io"

	core "rosetta-archive/internal"
)

// ArchiveFooter is the fixed final archive structure. It has no magic, version,
// or CRC; integrity is covered by HeaderCRC32, DataCRC32, and DirectoryCRC32.
type ArchiveFooter struct {
	CentralDirectoryOffset uint64
	CentralDirectorySize   uint64
	EntryCount             uint32
	FooterSize             uint32
}

type Footer = ArchiveFooter

func EncodeFooter(w io.Writer, f ArchiveFooter) error {
	if w == nil {
		return core.ErrNilWriter
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
		return ArchiveFooter{}, core.ErrNilReader
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

func ValidateFooter(f ArchiveFooter) error {
	if f.FooterSize != FooterSize {
		return core.ErrInvalidSize
	}
	return nil
}
