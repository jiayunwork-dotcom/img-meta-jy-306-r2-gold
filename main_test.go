package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImageFiles_UppercaseExt(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "shot.PNG")
	if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := imageFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 image, got %d %v", len(got), got)
	}
	if got[0] != p {
		t.Fatalf("got %q want %q", got[0], p)
	}
}
