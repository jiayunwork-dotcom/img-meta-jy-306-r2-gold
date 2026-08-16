// Package meta 提取图片的几何与格式元数据（基于标准库 image 包，无需外部依赖）。
package meta

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"image/color"
	"os"
)

// Meta 是图片的元数据。
type Meta struct {
	Path        string
	Format      string
	Width       int
	Height      int
	MegaPixels  float64
	ColorModel  string
	BitDepth    int
}

// ErrUnsupported 表示文件不是可解码的图片或已损坏。
var ErrUnsupported = errors.New("unsupported or unreadable image")

// Extract 读取图片文件头，提取格式、尺寸、色彩模型等元数据。
func Extract(path string) (Meta, error) {
	f, err := os.Open(path)
	if err != nil {
		return Meta{}, err
	}
	defer f.Close()
	cfg, format, err := image.DecodeConfig(f)
	if err != nil {
		return Meta{}, fmt.Errorf("%w: %v", ErrUnsupported, err)
	}
	return Meta{
		Path:       path,
		Format:     format,
		Width:      cfg.Width,
		Height:     cfg.Height,
		MegaPixels: float64(cfg.Width*cfg.Height) / 1e6,
		ColorModel: colorModelName(cfg.ColorModel),
		BitDepth:   bitDepth(cfg.ColorModel),
	}, nil
}

func colorModelName(cm color.Model) string {
	switch cm.Convert(color.White).(type) {
	case color.Gray, *color.Gray:
		return "gray"
	case color.Gray16, *color.Gray16, color.RGBA64, *color.RGBA64, color.NRGBA64, *color.NRGBA64:
		return "deep"
	case color.CMYK, *color.CMYK:
		return "cmyk"
	case color.YCbCr, *color.YCbCr:
		return "ycbcr"
	case color.RGBA, *color.RGBA, color.NRGBA, *color.NRGBA:
		return "rgb"
	default:
		return "other"
	}
}

func bitDepth(cm color.Model) int {
	switch cm.Convert(color.White).(type) {
	case color.Gray16, *color.Gray16, color.RGBA64, *color.RGBA64, color.NRGBA64, *color.NRGBA64:
		return 16
	default:
		return 8
	}
}
