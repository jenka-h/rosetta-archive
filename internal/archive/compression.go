package archive

import (
	"bytes"
	"fmt"
	"io"

	core "rosetta-archive/internal"
	"rosetta-archive/internal/format"
)

func encodePayload(method format.CompressionMethod, data []byte) ([]byte, error) {
	switch method {
	case format.CompressionStore:
		return data, nil
	case format.CompressionROSA1:
		return encodeRLE(data), nil
	default:
		return nil, fmt.Errorf("unsupported compression method %d", method)
	}
}

func decodePayload(method format.CompressionMethod, data []byte, uncompressedSize uint64) ([]byte, error) {
	switch method {
	case format.CompressionStore:
		if uint64(len(data)) != uncompressedSize {
			return nil, fmt.Errorf("stored payload size mismatch")
		}
		return data, nil
	case format.CompressionROSA1:
		return decodeRLE(data, uncompressedSize)
	default:
		return nil, fmt.Errorf("unsupported compression method %d", method)
	}
}

func (r *Reader) readPayload(entry format.DirectoryEntry) ([]byte, error) {
	compressed := make([]byte, entry.CompressedSize)
	if _, err := r.r.ReadAt(compressed, int64(entry.DataOffset)); err != nil {
		return nil, fmt.Errorf("read payload: %w", err)
	}
	decoded, err := decodePayload(entry.CompressionMethod, compressed, entry.UncompressedSize)
	if err != nil {
		return nil, err
	}
	if core.CRC32(decoded) != entry.DataCRC32 {
		return nil, fmt.Errorf("crc32 mismatch")
	}
	return decoded, nil
}

func payloadReader(data []byte) io.Reader {
	return bytes.NewReader(data)
}

// encodeRLE stores runs as count,value pairs. count is 1..255.
func encodeRLE(data []byte) []byte {
	if len(data) == 0 {
		return nil
	}
	out := make([]byte, 0, len(data))
	for i := 0; i < len(data); {
		b := data[i]
		count := 1
		for i+count < len(data) && data[i+count] == b && count < 255 {
			count++
		}
		out = append(out, byte(count), b)
		i += count
	}
	return out
}

func decodeRLE(data []byte, want uint64) ([]byte, error) {
	if len(data)%2 != 0 {
		return nil, fmt.Errorf("invalid rle payload length")
	}
	out := make([]byte, 0, want)
	for i := 0; i < len(data); i += 2 {
		count := int(data[i])
		if count == 0 {
			return nil, fmt.Errorf("invalid rle zero run")
		}
		for j := 0; j < count; j++ {
			out = append(out, data[i+1])
		}
		if uint64(len(out)) > want {
			return nil, fmt.Errorf("rle payload exceeds expected size")
		}
	}
	if uint64(len(out)) != want {
		return nil, fmt.Errorf("rle payload size mismatch")
	}
	return out, nil
}
