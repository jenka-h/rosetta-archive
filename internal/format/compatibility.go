package format

import core "rosetta-archive/internal"

// Compatibility describes how a reader should treat a major/minor version pair.
type Compatibility struct {
	SupportedMajor bool
	SupportedMinor bool
	Major          uint16
	Minor          uint16
}

// CheckCompatibility validates ROSA format version compatibility.
//
// Compatibility rules:
//   - different major versions are rejected;
//   - same major with newer minor versions is accepted;
//   - older minor versions are accepted.
func CheckCompatibility(major uint16, minor uint16) (Compatibility, error) {
	if major != FormatMajor {
		return Compatibility{}, core.ErrUnsupportedVersion
	}
	return Compatibility{
		SupportedMajor: true,
		SupportedMinor: minor <= FormatMinor,
		Major:          major,
		Minor:          minor,
	}, nil
}
