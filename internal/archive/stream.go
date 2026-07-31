package archive

import "io"

// StreamOptions controls sequential archive reading for large archives.
type StreamOptions struct {
	// VerifyDataCRC verifies each file payload as it is streamed.
	VerifyDataCRC bool
}

// StreamEntry is the callback payload for sequential archive iteration.
type StreamEntry struct {
	Entry Entry
	// Data contains the file payload for regular files. Directory entries have nil Data.
	Data io.Reader
}

// Stream iterates entries sequentially without requiring random access to the whole archive.
//
// TODO: implement after local entry parsing exists. Planned behavior:
//   - read and validate the archive header;
//   - decode each local entry record in file order;
//   - stream file payloads with io.LimitedReader;
//   - optionally verify CRC while streaming;
//   - stop before the central directory/footer without buffering full payloads.
func Stream(r io.Reader, options StreamOptions, visit func(StreamEntry) error) error {
	if r == nil {
		return notImplemented("stream archive: nil reader")
	}
	if visit == nil {
		return notImplemented("stream archive: nil visit callback")
	}
	return notImplemented("stream archive")
}

// OpenEntry opens one file payload by archive path using random access.
//
// TODO: implement after NewReader stores central-directory offsets. It must validate
// the requested entry metadata before returning a bounded reader.
func (r *Reader) OpenEntry(archivePath string) (io.Reader, Entry, error) {
	if r == nil {
		return nil, Entry{}, notImplemented("open archive entry: nil reader")
	}
	if archivePath == "" {
		return nil, Entry{}, notImplemented("open archive entry: missing path")
	}
	return nil, Entry{}, notImplemented("open archive entry")
}
