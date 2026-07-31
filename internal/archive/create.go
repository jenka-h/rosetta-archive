package archive

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	core "rosetta-archive/internal"
	"rosetta-archive/internal/format"
)

type sourceEntry struct {
	archivePath string
	sourcePath  string
	entryType   format.EntryType
	size        uint64
	crc32       uint32
	modTime     int64
}

func Create(archivePath string, sources []string) error {
	if archivePath == "" {
		return fmt.Errorf("create archive: missing archive path")
	}
	if len(sources) == 0 {
		return fmt.Errorf("create archive: missing source paths")
	}

	entries, err := collectSources(sources)
	if err != nil {
		return fmt.Errorf("collect sources: %w", err)
	}

	tmpPath := archivePath + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("create temporary archive: %w", err)
	}
	defer closeIfNeeded(out)

	if err := format.EncodeHeader(out, format.NewArchiveHeader()); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	offset := uint64(format.ArchiveHeaderSize)
	directory := make([]format.DirectoryEntry, 0, len(entries))

	for _, entry := range entries {
		dirEntry, nextOffset, err := writeEntry(out, entry, offset)
		if err != nil {
			_ = os.Remove(tmpPath)
			return err
		}
		directory = append(directory, dirEntry)
		offset = nextOffset
	}

	directoryOffset := offset
	directoryBytes, err := encodeDirectory(directory)
	if err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("encode central directory: %w", err)
	}
	if _, err := out.Write(directoryBytes); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("write central directory: %w", err)
	}

	if err := format.EncodeFooter(out, format.ArchiveFooter{
		CentralDirectoryOffset: directoryOffset,
		CentralDirectorySize:   uint64(len(directoryBytes)),
		EntryCount:             uint32(len(directory)),
	}); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}

	if err := out.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("close temporary archive: %w", err)
	}
	if err := os.Rename(tmpPath, archivePath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("replace archive: %w", err)
	}
	return nil
}

func writeEntry(out io.Writer, entry sourceEntry, offset uint64) (format.DirectoryEntry, uint64, error) {
	pathLen := uint64(len(entry.archivePath))
	header := format.EntryHeader{
		EntryType:           entry.entryType,
		CompressionMethod:   format.CompressionStore,
		PathLength:          uint32(pathLen),
		ExtraMetadataLength: 0,
		UncompressedSize:    entry.size,
		CompressedSize:      entry.size,
		ModificationTime:    entry.modTime,
		DataCRC32:           entry.crc32,
	}
	if entry.entryType == format.EntryTypeDirectory {
		header.UncompressedSize = 0
		header.CompressedSize = 0
		header.DataCRC32 = 0
	}
	if err := format.EncodeEntry(out, format.Entry{Header: header, Path: entry.archivePath}); err != nil {
		return format.DirectoryEntry{}, offset, fmt.Errorf("write entry %q: %w", entry.archivePath, err)
	}

	dataOffset := offset + uint64(format.EntryHeaderSize) + pathLen
	nextOffset := dataOffset
	if entry.entryType == format.EntryTypeFile {
		if err := copyFile(out, entry.sourcePath); err != nil {
			return format.DirectoryEntry{}, offset, fmt.Errorf("write data %q: %w", entry.archivePath, err)
		}
		nextOffset += entry.size
	}

	return format.DirectoryEntry{
		EntryType:         entry.entryType,
		CompressionMethod: format.CompressionStore,
		FileRecordOffset:  offset,
		DataOffset:        dataOffset,
		UncompressedSize:  header.UncompressedSize,
		CompressedSize:    header.CompressedSize,
		DataCRC32:         header.DataCRC32,
		Path:              entry.archivePath,
	}, nextOffset, nil
}

func collectSources(sources []string) ([]sourceEntry, error) {
	entries := make([]sourceEntry, 0)
	seen := make(map[string]struct{})
	for _, src := range sources {
		root := filepath.Clean(src)
		parent := filepath.Dir(root)
		if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			entry, err := sourceFromPath(parent, path, d)
			if err != nil {
				return err
			}
			if _, ok := seen[entry.archivePath]; ok {
				return fmt.Errorf("duplicate archive path %q", entry.archivePath)
			}
			seen[entry.archivePath] = struct{}{}
			entries = append(entries, entry)
			return nil
		}); err != nil {
			return nil, err
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].archivePath < entries[j].archivePath })
	return entries, nil
}

func sourceFromPath(parent string, path string, d fs.DirEntry) (sourceEntry, error) {
	info, err := d.Info()
	if err != nil {
		return sourceEntry{}, err
	}
	mode := info.Mode()
	if mode&os.ModeType != 0 && !mode.IsDir() {
		return sourceEntry{}, fmt.Errorf("unsupported filesystem object %q", path)
	}
	rel, err := filepath.Rel(parent, path)
	if err != nil {
		return sourceEntry{}, err
	}
	archivePath, err := core.NormalizeArchivePath(filepath.ToSlash(rel))
	if err != nil {
		return sourceEntry{}, fmt.Errorf("normalize %q: %w", path, err)
	}
	entry := sourceEntry{archivePath: archivePath, sourcePath: path, modTime: info.ModTime().Unix()}
	if info.IsDir() {
		entry.entryType = format.EntryTypeDirectory
		return entry, nil
	}
	if !info.Mode().IsRegular() {
		return sourceEntry{}, fmt.Errorf("unsupported filesystem object %q", path)
	}
	entry.entryType = format.EntryTypeFile
	entry.size = uint64(info.Size())
	entry.crc32, err = checksumFile(path)
	return entry, err
}

func checksumFile(path string) (uint32, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("open file for crc: %w", err)
	}
	defer closeIfNeeded(f)
	_, sum, err := core.CopyCRC32(io.Discard, f)
	if err != nil {
		return 0, fmt.Errorf("read file for crc: %w", err)
	}
	return sum, nil
}

func copyFile(w io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file data: %w", err)
	}
	defer closeIfNeeded(f)
	_, err = io.Copy(w, f)
	return err
}

func encodeDirectory(entries []format.DirectoryEntry) ([]byte, error) {
	var buf bytes.Buffer
	if err := format.EncodeDirectoryHeader(&buf, format.DirectoryHeader{EntryCount: uint32(len(entries))}); err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if err := format.EncodeDirectoryEntry(&buf, entry); err != nil {
			return nil, err
		}
	}
	data := buf.Bytes()
	format.ByteOrder.PutUint32(data[4:8], format.DirectoryCRC32(data))
	return data, nil
}
