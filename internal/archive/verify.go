package archive

import (
	"fmt"
	"io"
)

// Verify validates archive structure and CRC32 checksums.
//
// TODO: reuse Open and perform payload checksum verification after parsing exists.
func Verify(archivePath string) error {
	if archivePath == "" {
		return notImplemented("verify archive: missing archive path")
	}
	return notImplemented("verify archive")
}

// Info returns a high-level archive summary.
//
// TODO: derive this from the shared Reader once parsing exists.
func Info(archivePath string) (InfoSummary, error) {
	if archivePath == "" {
		return InfoSummary{}, notImplemented("archive info: missing archive path")
	}
	return InfoSummary{}, notImplemented("archive info")
}

// Inspect writes human-readable structural details about an archive.
//
// TODO: show header, metadata boundaries, payload ranges, footer, and checksum status.
func Inspect(archivePath string, out io.Writer) error {
	if archivePath == "" {
		return notImplemented("inspect archive: missing archive path")
	}
	if out == nil {
		return notImplemented("inspect archive: nil writer")
	}
	_, _ = fmt.Fprintln(out, "ROSA archive inspection is not implemented yet.")
	return notImplemented("inspect archive")
}
