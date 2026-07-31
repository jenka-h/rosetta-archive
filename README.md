# rosetta-archive

Implementation base for ROSA, a custom binary archive format written in Go.

This repository currently contains a project scaffold, archive operation stubs, and the first binary-format codec baseline in `internal/format`.

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

Archive operations currently return clear `not implemented` errors until create/list/extract/verify are wired to the format codecs.

## Development

Use only the Go standard library unless an external dependency is explicitly approved.

```sh
go test ./...
go vet ./...
```

## Current status

- Binary-format codec baseline in `internal/format`
- Archive operation stubs in `internal/archive`
- CLI command routing in `cmd/rosetta`
- Path normalization and CRC32 helper base in `internal`
- Format summary in `docs/specification.md`
