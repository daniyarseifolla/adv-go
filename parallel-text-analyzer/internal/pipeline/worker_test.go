package pipeline

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	content := "hello world\nfoo bar baz\n"
	os.WriteFile(path, []byte(content), 0644)

	stats, err := ProcessFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.Lines != 2 {
		t.Errorf("expected 2 lines, got %d", stats.Lines)
	}
	if stats.Words != 5 {
		t.Errorf("expected 5 words, got %d", stats.Words)
	}
	if stats.WordFreq["hello"] != 1 {
		t.Errorf("expected hello=1, got %d", stats.WordFreq["hello"])
	}
}

func TestProcessFileEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")
	os.WriteFile(path, []byte(""), 0644)

	stats, err := ProcessFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.Lines != 0 {
		t.Errorf("expected 0 lines, got %d", stats.Lines)
	}
	if stats.Words != 0 {
		t.Errorf("expected 0 words, got %d", stats.Words)
	}
}

func TestProcessFileNotFound(t *testing.T) {
	_, err := ProcessFile("/nonexistent/file.txt")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestWalkDir(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("world"), 0644)
	os.WriteFile(filepath.Join(dir, "c.log"), []byte("skip"), 0644)

	files, err := WalkDir(dir, ".txt", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}

func TestWalkDirSizeFilter(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "small.txt"), []byte("hi"), 0644)
	os.WriteFile(filepath.Join(dir, "big.txt"), []byte("hello world this is a bigger file with more content"), 0644)

	files, err := WalkDir(dir, ".txt", 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 file with min-size=10, got %d", len(files))
	}
}
