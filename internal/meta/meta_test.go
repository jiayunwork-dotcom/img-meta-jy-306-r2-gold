package meta

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func writePNG(t *testing.T, w, h int) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "img-*.png")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	img.Set(0, 0, color.RGBA{1, 2, 3, 255})
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestExtract_OK(t *testing.T) {
	p := writePNG(t, 1920, 1080)
	m, err := Extract(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Format != "png" {
		t.Fatalf("expected png, got %q", m.Format)
	}
	if m.Width != 1920 || m.Height != 1080 {
		t.Fatalf("bad dimensions: %dx%d", m.Width, m.Height)
	}
}

func TestExtract_Unsupported(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "x.txt")
	if err := os.WriteFile(p, []byte("not an image"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Extract(p); !errors.Is(err, ErrUnsupported) {
		t.Fatalf("expected ErrUnsupported, got %v", err)
	}
}

func TestClassify_ByAspect(t *testing.T) {
	land := Classify(Meta{Width: 1920, Height: 1080})
	if land.Aspect != "landscape" || land.ResTier != "fullhd" {
		t.Fatalf("bad landscape class: %+v", land)
	}
	port := Classify(Meta{Width: 600, Height: 800})
	if port.Aspect != "portrait" {
		t.Fatalf("bad portrait class: %+v", port)
	}
	sq := Classify(Meta{Width: 512, Height: 512})
	if sq.Aspect != "square" {
		t.Fatalf("bad square class: %+v", sq)
	}
}

func TestClassify_Square(t *testing.T) {
	got := Classify(Meta{Width: 512, Height: 512})
	if got.Aspect != "square" {
		t.Fatalf("equal width/height: aspect=%q", got.Aspect)
	}
}
