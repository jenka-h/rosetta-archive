# ROSA Manual Test Layout

Use this directory for local archive testing.

```text
tests/
├── in/         source files to archive
├── out/        generated .rosa files
└── extract/    extracted output
```

## Create archive

From `Zed/SISTER/rosetta-archive`:

```sh
go run ./cmd/rosetta create tests/out/example.rosa tests/in
```

## List archive

```sh
go run ./cmd/rosetta list tests/out/example.rosa
```

## Info

```sh
go run ./cmd/rosetta info tests/out/example.rosa
```

## Inspect

```sh
go run ./cmd/rosetta inspect tests/out/example.rosa
```

## Verify

```sh
go run ./cmd/rosetta verify tests/out/example.rosa
```

## Extract back

```sh
rm -rf tests/extract/*
go run ./cmd/rosetta extract tests/out/example.rosa tests/extract
```

The extracted files will appear under a folder matching the archived input directory, for example:

```text
tests/extract/in/
```

Compare with:

```sh
diff -r tests/in tests/extract/in
```
