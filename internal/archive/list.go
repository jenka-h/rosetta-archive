package archive

// List returns archive entries using the shared parser and validation path.
//
// TODO: call Open and return parsed entries once archive parsing exists.
func List(archivePath string) ([]Entry, error) {
	if archivePath == "" {
		return nil, notImplemented("list archive: missing archive path")
	}
	return nil, notImplemented("list archive")
}
