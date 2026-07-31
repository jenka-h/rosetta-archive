package archive

import (
	"fmt"
	"io"

	core "rosetta-archive/internal"
	"rosetta-archive/internal/format"
)

func Verify(archivePath string) error {
	if archivePath == "" {
		return fmt.Errorf("verify archive: missing archive path")
	}
	r, err := Open(archivePath)
	if err != nil {
		return err
	}
	defer closeIfNeeded(r)
	for _, entry := range r.directory {
		if entry.EntryType == format.EntryTypeFile {
			if err := r.verifyPayload(entry); err != nil {
				return fmt.Errorf("verify %q: %w", entry.Path, err)
			}
		}
	}
	return nil
}

func Info(archivePath string) (InfoSummary, error) {
	if archivePath == "" {
		return InfoSummary{}, fmt.Errorf("archive info: missing archive path")
	}
	r, err := Open(archivePath)
	if err != nil {
		return InfoSummary{}, err
	}
	defer closeIfNeeded(r)
	return r.summary(), nil
}

func Inspect(archivePath string, out io.Writer) error {
	if archivePath == "" {
		return fmt.Errorf("inspect archive: missing archive path")
	}
	if out == nil {
		return fmt.Errorf("inspect archive: nil writer")
	}
	r, err := Open(archivePath)
	if err != nil {
		return err
	}
	defer closeIfNeeded(r)

	printSummary(out, r.summary())
	fmt.Fprintln(out)
	fmt.Fprintln(out, "entries:")
	for i, entry := range r.directory {
		fmt.Fprintf(out, "  [%d]\n", i)
		fmt.Fprintf(out, "    path: %s\n", entry.Path)
		fmt.Fprintf(out, "    type: %s\n", entryTypeName(entry.EntryType))
		fmt.Fprintf(out, "    compression: %s\n", compressionName(entry.CompressionMethod))
		fmt.Fprintf(out, "    path_length: %d\n", entry.PathLength)
		fmt.Fprintf(out, "    file_record_offset: %d\n", entry.FileRecordOffset)
		fmt.Fprintf(out, "    data_offset: %d\n", entry.DataOffset)
		fmt.Fprintf(out, "    uncompressed_size: %d\n", entry.UncompressedSize)
		fmt.Fprintf(out, "    stored_size: %d\n", entry.CompressedSize)
		fmt.Fprintf(out, "    data_crc32: %08x\n", entry.DataCRC32)
	}
	return nil
}

func (r *Reader) summary() InfoSummary {
	s := InfoSummary{
		FormatName:   format.ArchiveMagic,
		MajorVersion: uint8(r.header.MajorVersion),
		MinorVersion: uint8(r.header.MinorVersion),
		ArchiveSize:  uint64(r.size),
		EntryCount:   uint32(len(r.directory)),
	}
	for _, entry := range r.directory {
		s.StoredSize += entry.CompressedSize
		s.UncompressedSize += entry.UncompressedSize
		if entry.EntryType == format.EntryTypeDirectory {
			s.DirectoryCount++
		} else {
			s.FileCount++
		}
	}
	return s
}

func printSummary(out io.Writer, s InfoSummary) {
	fmt.Fprintf(out, "format: %s\n", s.FormatName)
	fmt.Fprintf(out, "version: %d.%d\n", s.MajorVersion, s.MinorVersion)
	fmt.Fprintf(out, "archive_size: %d\n", s.ArchiveSize)
	fmt.Fprintf(out, "entry_count: %d\n", s.EntryCount)
	fmt.Fprintf(out, "file_count: %d\n", s.FileCount)
	fmt.Fprintf(out, "directory_count: %d\n", s.DirectoryCount)
	fmt.Fprintf(out, "stored_size: %d\n", s.StoredSize)
	fmt.Fprintf(out, "uncompressed_size: %d\n", s.UncompressedSize)
}

func entryTypeName(t format.EntryType) string {
	if t == format.EntryTypeDirectory {
		return "directory"
	}
	return "file"
}

func compressionName(m format.CompressionMethod) string {
	if m == format.CompressionROSA1 {
		return "rosa1"
	}
	return "store"
}

func (r *Reader) verifyPayload(entry format.DirectoryEntry) error {
	written, sum, err := core.CopyCRC32(io.Discard, io.NewSectionReader(r.r, int64(entry.DataOffset), int64(entry.CompressedSize)))
	if err != nil {
		return fmt.Errorf("read payload: %w", err)
	}
	if uint64(written) != entry.CompressedSize {
		return fmt.Errorf("payload size mismatch")
	}
	if sum != entry.DataCRC32 {
		return fmt.Errorf("crc32 mismatch")
	}
	return nil
}
