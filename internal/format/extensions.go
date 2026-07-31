package format

import (
	"io"

	core "rosetta-archive/internal"
)

// FeatureID identifies optional future storage features. Version 1 reserves these
// hooks but does not implement the feature payload formats.
type FeatureID uint16

const (
	FeatureCompression FeatureID = iota + 1
	FeatureEncryption
	FeatureDigitalSignature
)

// ExtensionRecord is a future TLV-style metadata record placeholder.
type ExtensionRecord struct {
	Type  FeatureID
	Value []byte
}

// EncodeExtensionRecord writes a future extension metadata record.
//
// TODO: implement only after the extension TLV layout is added to the written spec.
func EncodeExtensionRecord(w io.Writer, record ExtensionRecord) error {
	if w == nil {
		return core.ErrNilWriter
	}
	return core.ErrFeatureNotImplemented
}

// DecodeExtensionRecord reads a future extension metadata record.
//
// TODO: enforce length bounds before allocation when the extension layout is defined.
func DecodeExtensionRecord(r io.Reader) (ExtensionRecord, error) {
	if r == nil {
		return ExtensionRecord{}, core.ErrNilReader
	}
	return ExtensionRecord{}, core.ErrFeatureNotImplemented
}
