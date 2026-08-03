package archive

import "fmt"

func List(archivePath string) ([]Entry, error) {
	if archivePath == "" {
		return nil, fmt.Errorf("list archive: missing archive path")
	}
	r, err := Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer closeIfNeeded(r)
	return r.Entries(), nil
}
