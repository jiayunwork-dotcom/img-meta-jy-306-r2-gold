//go:build ignore

// 该文件用于生成示例图片，不参与正常构建（go build ./... 会跳过）。
// 运行：go run example/gen.go
package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func make(dir, name string, w, h int) {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), 128, 255})
		}
	}
	f, err := os.Create(dir + "/" + name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func main() {
	dir := "example"
	make(dir, "landscape.png", 1920, 1080)
	make(dir, "portrait.png", 600, 800)
	make(dir, "square.png", 512, 512)
}
