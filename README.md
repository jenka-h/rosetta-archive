# rosetta-archive

`rosetta-archive` is a Go implementation of ROSA, a custom binary archive format.

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
go run ./cmd/rosetta create tests/out/example.rosa tests/in
```

### Create With Compression

Uses `ROSA1` RLE compression.

```sh
go run ./cmd/rosetta create-compressed tests/out/example.rosa tests/in
```

### List

```sh
go run ./cmd/rosetta list tests/out/example.rosa
```

### Info

```sh
go run ./cmd/rosetta info tests/out/example.rosa
```

### Inspect

```sh
go run ./cmd/rosetta inspect tests/out/example.rosa
```

### Verify

```sh
go run ./cmd/rosetta verify tests/out/example.rosa
```

### Extract

```sh
rm -rf tests/extract/*
go run ./cmd/rosetta extract tests/out/example.rosa tests/extract
```

### Read One File With Random Access

```sh
go run ./cmd/rosetta cat tests/out/example.rosa in/file.txt
```

### Recover Damaged Archive Entries

```sh
go run ./cmd/rosetta recover tests/out/example.rosa
```

## Test Commands

```sh
go test ./...
go vet ./...
```

From `/home/real/Zed`:

```sh
go -C SISTER/rosetta-archive test ./...
go -C SISTER/rosetta-archive vet ./...
```
