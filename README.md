# rosetta-archive

Implementation base for ROSA, a custom binary archive format written in Go.

This repository currently contains only the project scaffold and operation stubs. The archive header, metadata records, payload layout, footer, and checksum coverage are intentionally not defined yet.

## Commands

```sh
go run ./cmd/rosetta --help
```

Planned commands:

- `rosetta create <archive.rosa> <path> [path...]`
- `rosetta list <archive.rosa>`
- `rosetta info <archive.rosa>`
- `rosetta inspect <archive.rosa>`
- `rosetta verify <archive.rosa>`
- `rosetta extract <archive.rosa> <destination>`

Each operation currently returns a clear `not implemented` error until the binary format is specified.

## Development

Use only the Go standard library unless an external dependency is explicitly approved.

```sh
go test ./...
go vet ./...
```

## Current status

- Temporary empty format package in `internal/format`
- Archive operation stubs in `internal/archive`
- CLI command routing in `cmd/rosetta`
- Path normalization and CRC32 helper base in `internal`
- Placeholder specification in `docs/specification.md`
