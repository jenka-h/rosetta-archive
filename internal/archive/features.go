package archive

import (
	"fmt"

	"rosetta-archive/internal/format"
)

func CreateWithCompression(archivePath string, sources []string, method format.CompressionMethod) error {
	switch method {
	case format.CompressionStore, format.CompressionROSA1:
		return createArchive(archivePath, sources, method)
	default:
		return fmt.Errorf("create archive with compression: unsupported method %d", method)
	}
}
