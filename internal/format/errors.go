package format

import "errors"

var (
	// ErrFormatNotDefined is returned until the ROSA binary format layout is specified.
	ErrFormatNotDefined = errors.New("rosa binary format is not defined yet")
	ErrNilReader        = errors.New("nil reader")
	ErrNilWriter        = errors.New("nil writer")
)
