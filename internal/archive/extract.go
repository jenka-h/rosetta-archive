package archive

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	core "rosetta-archive/internal"
	"rosetta-archive/internal/format"
)

func Extract(archivePath string, destination string) error {
	if archivePath == "" {
		return fmt.Errorf("extract archive: missing archive path")
	}
	if destination == "" {
		return fmt.Errorf("extract archive: missing destination")
	}
	r, err := Open(archivePath)
	if err != nil {
		return err
	}
	defer closeIfNeeded(r)

	root, err := filepath.Abs(destination)
	if err != nil {
		return fmt.Errorf("resolve destination: %w", err)
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return fmt.Errorf("create destination: %w", err)
	}

	for _, entry := range r.directory {
		target, err := safeJoin(root, entry.Path)
		if err != nil {
			return fmt.Errorf("validate extract path %q: %w", entry.Path, err)
		}
		if entry.EntryType == format.EntryTypeDirectory {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("create directory %q: %w", target, err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("create parent directory %q: %w", filepath.Dir(target), err)
		}
		if err := r.extractFile(entry, target); err != nil {
			return err
		}
	}
	return nil
}

func safeJoin(root string, archivePath string) (string, error) {
	clean, err := core.NormalizeArchivePath(archivePath)
	if err != nil {
		return "", err
	}
	target, err := filepath.Abs(filepath.Join(root, filepath.FromSlash(clean)))
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || filepath.IsAbs(rel) {
		return "", fmt.Errorf("path escapes destination")
	}
	return target, nil
}

func (r *Reader) extractFile(entry format.DirectoryEntry, target string) error {
	tmp, err := os.CreateTemp(filepath.Dir(target), ".rosetta-*")
	if err != nil {
		return fmt.Errorf("create temporary file for %q: %w", target, err)
	}
	tmpPath := tmp.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	written, sum, err := core.CopyCRC32(tmp, io.NewSectionReader(r.r, int64(entry.DataOffset), int64(entry.CompressedSize)))
	if closeErr := tmp.Close(); err == nil && closeErr != nil {
		err = closeErr
	}
	if err != nil {
		return fmt.Errorf("extract file %q: %w", entry.Path, err)
	}
	if uint64(written) != entry.UncompressedSize {
		return fmt.Errorf("extract file %q: size mismatch", entry.Path)
	}
	if sum != entry.DataCRC32 {
		return fmt.Errorf("extract file %q: crc32 mismatch", entry.Path)
	}
	if err := os.Rename(tmpPath, target); err != nil {
		return fmt.Errorf("replace extracted file %q: %w", target, err)
	}
	return nil
}
