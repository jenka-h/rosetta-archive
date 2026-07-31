package format

import (
	"io"

	core "rosetta-archive/internal"
)

// FeatureID identifies optional extension features. ROSA currently reserves only
// compression as an extension hook.
type FeatureID uint16

const FeatureCompression FeatureID = 1

// ExtensionRecord is a future compression metadata placeholder.
type ExtensionRecord struct {
	Type  FeatureID
	Value []byte
}

// EncodeExtensionRecord writes a future compression extension metadata record.
//
// TODO: implement only after the compression extension TLV layout is added to the written spec.
func EncodeExtensionRecord(w io.Writer, record ExtensionRecord) error {
	if w == nil {
		return core.ErrNilWriter
	}
	return core.ErrFeatureNotImplemented
}

// DecodeExtensionRecord reads a future compression extension metadata record.
//
// TODO: enforce length bounds before allocation when the compression extension layout is defined.
func DecodeExtensionRecord(r io.Reader) (ExtensionRecord, error) {
	if r == nil {
		return ExtensionRecord{}, core.ErrNilReader
	}
	return ExtensionRecord{}, core.ErrFeatureNotImplemented
}
