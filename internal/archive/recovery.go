package archive

import "io"

// RecoveryOptions controls best-effort recovery from a partially damaged archive.
//
// Recovery is a bonus feature and is intentionally only a stub until the required
// create/list/extract/verify path is complete.
type RecoveryOptions struct {
	// AllowPartialEntries permits returning entries whose metadata can be parsed even
	// when unrelated archive regions are corrupt. File payloads must still be verified
	// before extraction.
	AllowPartialEntries bool
	// ScanForEntryRecords permits sequential scanning for plausible local entry
	// records when the central directory or footer is damaged.
	ScanForEntryRecords bool
}

// RecoveryReport describes entries and corruption found during best-effort recovery.
type RecoveryReport struct {
	RecoveredEntries []Entry
	Warnings         []string
}

// Recover opens an archive from disk and attempts best-effort recovery.
//
// TODO: implement after normal parsing/verification works. Planned behavior:
//   - read valid header when present;
//   - locate footer and central directory when intact;
//   - optionally scan for structurally valid entry records when directory/footer is corrupt;
//   - validate every recovered path, size, offset, and CRC before exposing it;
//   - never extract unverified bytes silently.
func Recover(archivePath string, options RecoveryOptions) (RecoveryReport, error) {
	if archivePath == "" {
		return RecoveryReport{}, notImplemented("recover archive: missing archive path")
	}
	return RecoveryReport{}, notImplemented("recover archive")
}

// RecoverReader attempts best-effort recovery from an existing random-access reader.
//
// TODO: share implementation with Recover and keep all bounds checks centralized.
func RecoverReader(r io.ReaderAt, size int64, options RecoveryOptions) (RecoveryReport, error) {
	if r == nil {
		return RecoveryReport{}, notImplemented("recover archive reader: nil reader")
	}
	if size < 0 {
		return RecoveryReport{}, notImplemented("recover archive reader: negative size")
	}
	return RecoveryReport{}, notImplemented("recover archive reader")
}
