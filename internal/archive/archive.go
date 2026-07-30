package archive

import (
	"errors"
	"fmt"
)

// ErrNotImplemented is returned by archive operations until the ROSA v1 format is specified.
var ErrNotImplemented = errors.New("rosa archive operation not implemented")

// EntryType identifies the kind of filesystem object stored in an archive.
type EntryType uint8

const (
	EntryTypeFile EntryType = iota + 1
	EntryTypeDirectory
)

// Entry describes one archive member after parsing metadata.
type Entry struct {
	Path   string
	Type   EntryType
	Size   uint64
	CRC32  uint32
	Offset uint64
}

// InfoSummary contains high-level archive metadata.
type InfoSummary struct {
	EntryCount uint64
}

func notImplemented(operation string) error {
	return fmt.Errorf("%s: %w", operation, ErrNotImplemented)
}
