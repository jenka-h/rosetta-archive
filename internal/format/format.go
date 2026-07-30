package format

const (
	// VersionMajor is intentionally zero until the ROSA binary format is specified.
	VersionMajor uint16 = 0
	// VersionMinor is intentionally zero until the ROSA binary format is specified.
	VersionMinor uint16 = 0
)

// ByteOrder is intentionally not defined yet. Use encoding/binary.LittleEndian
// explicitly when the format fields are specified.
