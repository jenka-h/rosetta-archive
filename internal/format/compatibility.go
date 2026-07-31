package format

import core "rosetta-archive/internal"

func CheckCompatibility(major uint16, minor uint16) error {
	if major != FormatMajor {
		return core.ErrUnsupportedVersion
	}
	return nil
}

func IsKnownMinorVersion(minor uint16) bool {
	return minor <= FormatMinor
}
