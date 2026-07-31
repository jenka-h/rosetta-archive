package archive

// CreateOptions holds optional archive creation features. Required ROSA v1 behavior
// should be implemented before compression is enabled.
type CreateOptions struct {
	Compress bool
}

// CreateWithOptions creates an archive with optional future compression.
//
// TODO: reject compression until it has a written binary spec, tests, and compatibility rules.
func CreateWithOptions(archivePath string, sources []string, options CreateOptions) error {
	if archivePath == "" {
		return notImplemented("create archive with options: missing archive path")
	}
	if len(sources) == 0 {
		return notImplemented("create archive with options: missing source paths")
	}
	if options.Compress {
		return notImplemented("create archive with options: compression")
	}
	return Create(archivePath, sources)
}
