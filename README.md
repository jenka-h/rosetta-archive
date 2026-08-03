# rosetta-archive

`rosetta-archive` is a Go implementation of ROSA, a custom binary archive format.

## Video
[Video Demo](https://youtu.be/uKs6tdSyooQ)

## Project Structure

```text
cmd/rosetta/          CLI entrypoint and command parsing
internal/archive/     Archive create/open/list/info/inspect/verify/extract logic
internal/format/      Binary format structs, validation, marshal/unmarshal codecs
internal/             Shared checksum, path validation, and errors
docs/specification.md Binary format summary
tests/                Manual test input/output/extract folders
```

## Commands

Run commands from the repository root:

```sh
go run ./cmd/rosetta --help
```

### Create

```sh
go run ./cmd/rosetta create <archive_path> <input_path> [input_path...]
```

### Create With Compression

Uses `ROSA1` RLE compression.

```sh
go run ./cmd/rosetta create-compressed <archive_path> <input_path> [input_path...]
```

### List

```sh
go run ./cmd/rosetta list <archive_path>
```

### Info

```sh
go run ./cmd/rosetta info <archive_path>
```

### Inspect

```sh
go run ./cmd/rosetta inspect <archive_path>
```

### Verify

```sh
go run ./cmd/rosetta verify <archive_path>
```

### Extract

```sh
go run ./cmd/rosetta extract <archive_path> <destination_path>
```

### Read One File With Random Access

```sh
go run ./cmd/rosetta cat <archive_path> <entry_path>
```

### Recover Damaged Archive Entries

```sh
go run ./cmd/rosetta recover <archive_path>
```
