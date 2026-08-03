package internal

import (
	"hash/crc32"
	"io"
)

func CRC32(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

func CopyCRC32(dst io.Writer, src io.Reader) (written int64, sum uint32, err error) {
	h := crc32.NewIEEE()
	written, err = io.Copy(dst, io.TeeReader(src, h))
	return written, h.Sum32(), err
}
