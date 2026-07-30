package archive

import (
	"errors"
	"fmt"
)

var ErrNotImplemented = errors.New("rosa archive operation not implemented")

type EntryType uint8

// membedakan directory dengan file biasa, bisa konsider symbolic link, tapi struktur filesystem symbolic jarang
// dipakai, sehingga untuk sementara entry hanya dua berikut ini
const (
	EntryTypeFile EntryType = iota + 1
	EntryTypeDirectory
)

// entry sementara
type Entry struct {
	Path   string
	Type   EntryType
	Size   uint64
	CRC32  uint32
	Offset uint64
}

// metadata sementara
type InfoSummary struct {
	EntryCount uint64
}

func notImplemented(operation string) error {
	return fmt.Errorf("%s: %w", operation, ErrNotImplemented)
}
