package archive

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"rosetta-archive/internal/format"
)

func TestRoundTripCreateVerifyExtractInfoInspectAndStream(t *testing.T) {
	root := t.TempDir()
	input := filepath.Join(root, "in")
	if err := os.MkdirAll(filepath.Join(input, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(input, "hello.txt"), []byte("hello rosa\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	binary := []byte{0, 1, 2, 3, 127, 128, 254, 255}
	if err := os.WriteFile(filepath.Join(input, "nested", "sample.bin"), binary, 0o644); err != nil {
		t.Fatal(err)
	}

	archivePath := filepath.Join(root, "out.rosa")
	if err := Create(archivePath, []string{input}); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if err := Verify(archivePath); err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}

	info, err := Info(archivePath)
	if err != nil {
		t.Fatalf("Info returned error: %v", err)
	}
	if info.FormatName != "ROSA" || info.FileCount != 2 || info.DirectoryCount != 2 || info.EntryCount != 4 {
		t.Fatalf("unexpected info summary: %+v", info)
	}
	if info.StoredSize != uint64(len("hello rosa\n")+len(binary)) || info.UncompressedSize != info.StoredSize {
		t.Fatalf("unexpected sizes: %+v", info)
	}

	var inspect bytes.Buffer
	if err := Inspect(archivePath, &inspect); err != nil {
		t.Fatalf("Inspect returned error: %v", err)
	}
	if !strings.Contains(inspect.String(), "hello.txt") || !strings.Contains(inspect.String(), "sample.bin") {
		t.Fatalf("inspect output missing file metadata:\n%s", inspect.String())
	}

	entries, err := List(archivePath)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(entries) != 4 {
		t.Fatalf("entry count = %d, want 4", len(entries))
	}

	opened, err := Open(archivePath)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer closeIfNeeded(opened)
	payload, _, err := opened.OpenEntry("in/nested/sample.bin")
	if err != nil {
		t.Fatalf("OpenEntry returned error: %v", err)
	}
	gotPayload, err := io.ReadAll(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotPayload, binary) {
		t.Fatalf("OpenEntry payload = %v, want %v", gotPayload, binary)
	}

	archiveBytes, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	seen := 0
	if err := Stream(bytes.NewReader(archiveBytes), func(entry Entry, data io.Reader) error {
		seen++
		if entry.Path == "in/hello.txt" {
			payload, err := io.ReadAll(data)
			if err != nil {
				return err
			}
			if string(payload) != "hello rosa\n" {
				t.Fatalf("stream payload = %q", payload)
			}
		}
		return nil
	}); err != nil {
		t.Fatalf("Stream returned error: %v", err)
	}
	if seen != 4 {
		t.Fatalf("stream visited %d entries, want 4", seen)
	}

	extractDir := filepath.Join(root, "extract")
	if err := Extract(archivePath, extractDir); err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}
	gotText, err := os.ReadFile(filepath.Join(extractDir, "in", "hello.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(gotText) != "hello rosa\n" {
		t.Fatalf("extracted text = %q", gotText)
	}
	gotBinary, err := os.ReadFile(filepath.Join(extractDir, "in", "nested", "sample.bin"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotBinary, binary) {
		t.Fatalf("extracted binary = %v, want %v", gotBinary, binary)
	}
}

func TestCompressionAndRecovery(t *testing.T) {
	root := t.TempDir()
	input := filepath.Join(root, "in")
	if err := os.MkdirAll(input, 0o755); err != nil {
		t.Fatal(err)
	}
	content := bytes.Repeat([]byte("aaaaabbbbbcccccccccc"), 32)
	if err := os.WriteFile(filepath.Join(input, "runs.txt"), content, 0o644); err != nil {
		t.Fatal(err)
	}

	archivePath := filepath.Join(root, "compressed.rosa")
	if err := CreateWithCompression(archivePath, []string{input}, format.CompressionROSA1); err != nil {
		t.Fatalf("CreateWithCompression returned error: %v", err)
	}
	if err := Verify(archivePath); err != nil {
		t.Fatalf("Verify compressed archive returned error: %v", err)
	}
	opened, err := Open(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	defer closeIfNeeded(opened)
	payload, _, err := opened.OpenEntry("in/runs.txt")
	if err != nil {
		t.Fatal(err)
	}
	got, err := io.ReadAll(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, content) {
		t.Fatal("compressed payload did not decode to original content")
	}

	corrupt, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < int(format.FooterSize); i++ {
		corrupt[len(corrupt)-1-i] = 0
	}
	recovered, err := RecoverReader(bytes.NewReader(corrupt), int64(len(corrupt)))
	if err != nil {
		t.Fatalf("RecoverReader returned error: %v", err)
	}
	if len(recovered) == 0 {
		t.Fatal("RecoverReader recovered no entries")
	}
}
