package format

func marshalHeader(h ArchiveHeader) []byte {
	buf := make([]byte, ArchiveHeaderSize)
	putMagic(buf[0:4], ArchiveMagic)
	ByteOrder.PutUint16(buf[4:6], h.MajorVersion)
	ByteOrder.PutUint16(buf[6:8], h.MinorVersion)
	ByteOrder.PutUint32(buf[8:12], h.HeaderSize)
	ByteOrder.PutUint32(buf[12:16], h.HeaderCRC32)
	return buf
}

func unmarshalHeader(buf []byte) ArchiveHeader {
	return ArchiveHeader{
		MajorVersion: ByteOrder.Uint16(buf[4:6]),
		MinorVersion: ByteOrder.Uint16(buf[6:8]),
		HeaderSize:   ByteOrder.Uint32(buf[8:12]),
		HeaderCRC32:  ByteOrder.Uint32(buf[12:16]),
	}
}

func marshalDirectoryHeader(h DirectoryHeader) []byte {
	buf := make([]byte, DirectoryHeaderSize)
	ByteOrder.PutUint32(buf[0:4], h.EntryCount)
	ByteOrder.PutUint32(buf[4:8], h.DirectoryCRC32)
	ByteOrder.PutUint32(buf[8:12], h.Reserved)
	return buf
}

func unmarshalDirectoryHeader(buf []byte) DirectoryHeader {
	return DirectoryHeader{
		EntryCount:     ByteOrder.Uint32(buf[0:4]),
		DirectoryCRC32: ByteOrder.Uint32(buf[4:8]),
		Reserved:       ByteOrder.Uint32(buf[8:12]),
	}
}

func marshalFooter(f ArchiveFooter) []byte {
	buf := make([]byte, FooterSize)
	ByteOrder.PutUint16(buf[0:2], f.MajorVersion)
	ByteOrder.PutUint16(buf[2:4], f.MinorVersion)
	ByteOrder.PutUint64(buf[4:12], f.CentralDirectoryOffset)
	ByteOrder.PutUint64(buf[12:20], f.CentralDirectorySize)
	ByteOrder.PutUint32(buf[20:24], f.EntryCount)
	ByteOrder.PutUint32(buf[24:28], f.FooterSize)
	return buf
}

func unmarshalFooter(buf []byte) ArchiveFooter {
	return ArchiveFooter{
		MajorVersion:           ByteOrder.Uint16(buf[0:2]),
		MinorVersion:           ByteOrder.Uint16(buf[2:4]),
		CentralDirectoryOffset: ByteOrder.Uint64(buf[4:12]),
		CentralDirectorySize:   ByteOrder.Uint64(buf[12:20]),
		EntryCount:             ByteOrder.Uint32(buf[20:24]),
		FooterSize:             ByteOrder.Uint32(buf[24:28]),
	}
}
