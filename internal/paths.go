package internal

import (
	"errors"
	"path"
	"strings"
	"unicode/utf8"
)

const (
	MaxEntries    = 1_000_000 // maximum entries di archive
	MaxPathLength = 4096      // maximum length dari path di archive
)

// TODO: extend this when traversal code defines exact source-root behavior.
func NormalizeArchivePath(p string) (string, error) {
	if p == "" {
		return "", errors.New("empty archive path")
	}
	if len(p) > MaxPathLength {
		return "", errors.New("archive path exceeds maximum length")
	}
	if !utf8.ValidString(p) {
		return "", errors.New("archive path is not valid UTF-8")
	}
	if strings.ContainsRune(p, '\x00') {
		return "", errors.New("archive path contains null byte")
	}
	p = strings.ReplaceAll(p, "\\", "/")
	if strings.HasPrefix(p, "/") {
		return "", errors.New("archive path must be relative")
	}
	if strings.Contains(p, "//") {
		return "", errors.New("archive path contains duplicate separators")
	}
	for _, part := range strings.Split(p, "/") {
		if part == "." || part == ".." {
			return "", errors.New("archive path contains unsafe component")
		}
	}
	clean := path.Clean(p)
	if clean == "." || strings.HasPrefix(clean, "../") || clean == ".." {
		return "", errors.New("archive path escapes root")
	}
	return clean, nil
}
