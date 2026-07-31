package archive

// CreateOptions holds optional archive creation features. Required ROSA v1 behavior
// should be implemented before any bonus feature is enabled.
type CreateOptions struct {
	Compress         bool
	Encrypt          bool
	DigitalSignature bool
}

// CreateWithOptions creates an archive with optional future storage features.
//
// TODO: reject unsupported options until compression, encryption, and digital
// signature each have a written binary spec, tests, and compatibility rules.
func CreateWithOptions(archivePath string, sources []string, options CreateOptions) error {
	if archivePath == "" {
		return notImplemented("create archive with options: missing archive path")
	}
	if len(sources) == 0 {
		return notImplemented("create archive with options: missing source paths")
	}
	if options.Compress || options.Encrypt || options.DigitalSignature {
		return notImplemented("create archive with options: requested bonus feature")
	}
	return Create(archivePath, sources)
}
