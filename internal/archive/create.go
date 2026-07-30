package archive

// Create writes a ROSA archive at archivePath from the provided source paths.
//
// TODO: implement deterministic traversal, metadata generation, payload writing,
// checksum calculation, and atomic archive replacement after the binary format is defined.
func Create(archivePath string, sources []string) error {
	if archivePath == "" {
		return notImplemented("create archive: missing archive path")
	}
	if len(sources) == 0 {
		return notImplemented("create archive: missing source paths")
	}
	return notImplemented("create archive")
}
