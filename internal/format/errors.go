package format

import "errors"

var (
	ErrNilReader                = errors.New("nil reader")
	ErrNilWriter                = errors.New("nil writer")
	ErrInvalidMagic             = errors.New("invalid magic")
	ErrUnsupportedVersion       = errors.New("unsupported version")
	ErrInvalidSize              = errors.New("invalid format size")
	ErrReservedNonZero          = errors.New("reserved field is non-zero")
	ErrUnknownFlags             = errors.New("unknown reserved flag bits are set")
	ErrInvalidEntryType         = errors.New("invalid entry type")
	ErrInvalidCompressionMethod = errors.New("invalid compression method")
	ErrInvalidEntrySizes        = errors.New("invalid entry sizes")
	ErrInvalidEntryCRC          = errors.New("invalid entry crc field")
	ErrInvalidModificationTime  = errors.New("invalid modification time field")
	ErrInvalidPath              = errors.New("invalid archive path")
	ErrPathTooLong              = errors.New("archive path exceeds maximum length")
	ErrExtraMetadataTooLong     = errors.New("extra metadata exceeds maximum length")
	ErrInvalidFooterCRC32       = errors.New("invalid footer crc32")
)
