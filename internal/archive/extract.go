package archive

// cmd untuk mengekstrak arsip ROSA ke direktori tujuan
//
// TODO: validate archive paths, prevent traversal, write temporary files first,
// verify size and CRC32 before rename, and reject unsupported entry types.
func Extract(archivePath string, destination string) error {
	if archivePath == "" {
		return notImplemented("extract archive: missing archive path")
	}
	if destination == "" {
		return notImplemented("extract archive: missing destination")
	}
	return notImplemented("extract archive")
}
