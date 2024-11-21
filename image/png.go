package image

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
)

// PNGCompression 将png图片压缩错8bit的color map模式
func PNGCompression(body []byte) ([]byte, error) {
	// 解码PNG图像
	img, err := png.Decode(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("failed to decode png: %v", err)
	}

	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

	if !bytes.HasPrefix(body, pngHeader) {
		return body, errors.New("The picture format is incorrect.")
	}

	if _, ok := img.(*image.Gray); ok {
		return body, nil
	}

	// 创建调色板
	palette := createPalette()

	// 创建一个新的调色板图像
	bounds := img.Bounds()
	palettedImg := image.NewPaletted(bounds, palette)

	// 将原始图像的像素映射到调色板图像
	draw.FloydSteinberg.Draw(palettedImg, bounds, img, image.Point{})

	// 创建输出文件
	var buf bytes.Buffer
	// 将调色板图像编码为PNG并保存
	err = png.Encode(&buf, palettedImg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// createPalette 创建一个简单的调色板，这里我们使用256种颜色
func createPalette() color.Palette {
	var palette color.Palette
	for r := 0; r < 256; r += 51 {
		for g := 0; g < 256; g += 51 {
			for b := 0; b < 256; b += 51 {
				palette = append(palette, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
			}
		}
	}
	return palette
}
