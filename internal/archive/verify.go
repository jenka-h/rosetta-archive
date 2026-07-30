package archive

import (
	"fmt"
	"io"
)

// cmd untuk memverifikasi struktur arsip ROSA dan checksum CRC32
//
// TODO: reuse Open and perform payload checksum verification after parsing exists.
func Verify(archivePath string) error {
	if archivePath == "" {
		return notImplemented("verify archive: missing archive path")
	}
	return notImplemented("verify archive")
}

// cmd untuk mendapatkan ringkasan informasi arsip ROSA
//
// TODO: derive this from the shared Reader once parsing exists.
func Info(archivePath string) (InfoSummary, error) {
	if archivePath == "" {
		return InfoSummary{}, notImplemented("archive info: missing archive path")
	}
	return InfoSummary{}, notImplemented("archive info")
}

// cmd untuk inspeksi detail struktur arsip ROSA
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
