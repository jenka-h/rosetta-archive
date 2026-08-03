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

func marshalEntryHeader(h EntryHeader) []byte {
	buf := make([]byte, EntryHeaderSize)
	buf[0] = byte(h.EntryType)
	buf[1] = byte(h.CompressionMethod)
	ByteOrder.PutUint32(buf[2:6], h.PathLength)
	ByteOrder.PutUint32(buf[6:10], h.ExtraMetadataLength)
	ByteOrder.PutUint64(buf[10:18], h.UncompressedSize)
	ByteOrder.PutUint64(buf[18:26], h.CompressedSize)
	ByteOrder.PutUint64(buf[26:34], uint64(h.ModificationTime))
	ByteOrder.PutUint32(buf[34:38], h.DataCRC32)
	ByteOrder.PutUint32(buf[38:42], h.Reserved)
	return buf
}

func unmarshalEntryHeader(buf []byte) EntryHeader {
	return EntryHeader{
		EntryType:           EntryType(buf[0]),
		CompressionMethod:   CompressionMethod(buf[1]),
		PathLength:          ByteOrder.Uint32(buf[2:6]),
		ExtraMetadataLength: ByteOrder.Uint32(buf[6:10]),
		UncompressedSize:    ByteOrder.Uint64(buf[10:18]),
		CompressedSize:      ByteOrder.Uint64(buf[18:26]),
		ModificationTime:    int64(ByteOrder.Uint64(buf[26:34])),
		DataCRC32:           ByteOrder.Uint32(buf[34:38]),
		Reserved:            ByteOrder.Uint32(buf[38:42]),
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

func marshalDirectoryEntry(e DirectoryEntry) []byte {
	buf := make([]byte, DirectoryEntrySize)
	buf[0] = byte(e.EntryType)
	buf[1] = byte(e.CompressionMethod)
	ByteOrder.PutUint32(buf[2:6], e.PathLength)
	ByteOrder.PutUint64(buf[6:14], e.FileRecordOffset)
	ByteOrder.PutUint64(buf[14:22], e.DataOffset)
	ByteOrder.PutUint64(buf[22:30], e.UncompressedSize)
	ByteOrder.PutUint64(buf[30:38], e.CompressedSize)
	ByteOrder.PutUint32(buf[38:42], e.DataCRC32)
	ByteOrder.PutUint32(buf[42:46], e.Reserved)
	return buf
}

func unmarshalDirectoryEntry(buf []byte) DirectoryEntry {
	return DirectoryEntry{
		EntryType:         EntryType(buf[0]),
		CompressionMethod: CompressionMethod(buf[1]),
		PathLength:        ByteOrder.Uint32(buf[2:6]),
		FileRecordOffset:  ByteOrder.Uint64(buf[6:14]),
		DataOffset:        ByteOrder.Uint64(buf[14:22]),
		UncompressedSize:  ByteOrder.Uint64(buf[22:30]),
		CompressedSize:    ByteOrder.Uint64(buf[30:38]),
		DataCRC32:         ByteOrder.Uint32(buf[38:42]),
		Reserved:          ByteOrder.Uint32(buf[42:46]),
	}
}

func marshalFooter(f ArchiveFooter) []byte {
	buf := make([]byte, FooterSize)
	ByteOrder.PutUint64(buf[0:8], f.CentralDirectoryOffset)
	ByteOrder.PutUint64(buf[8:16], f.CentralDirectorySize)
	ByteOrder.PutUint32(buf[16:20], f.EntryCount)
	ByteOrder.PutUint32(buf[20:24], f.FooterSize)
	return buf
}

func unmarshalFooter(buf []byte) ArchiveFooter {
	return ArchiveFooter{
		CentralDirectoryOffset: ByteOrder.Uint64(buf[0:8]),
		CentralDirectorySize:   ByteOrder.Uint64(buf[8:16]),
		EntryCount:             ByteOrder.Uint32(buf[16:20]),
		FooterSize:             ByteOrder.Uint32(buf[20:24]),
	}
}
