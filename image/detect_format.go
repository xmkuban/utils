package image

import "bytes"

func DetectImageFormat(data []byte) string {
	// JPEG (JPG) 文件头
	jpegHeader := []byte{0xFF, 0xD8, 0xFF}
	// PNG 文件头
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	// GIF 文件头
	gifHeader1 := []byte{0x47, 0x49, 0x46, 0x38, 0x37, 0x61}
	gifHeader2 := []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}
	// BMP 文件头
	bmpHeader := []byte{0x42, 0x4D}
	// TIFF 文件头
	tiffHeader1 := []byte{0x49, 0x49, 0x2A, 0x00}
	tiffHeader2 := []byte{0x4D, 0x4D, 0x00, 0x2A}
	// WebP 文件头
	webpHeader := []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50}
	// ICO 文件头
	icoHeader := []byte{0x00, 0x00, 0x01, 0x00}
	// SVG 文件头
	svgHeader := []byte{0x3C, 0x3F, 0x78, 0x6D, 0x6C}
	// HEIC 文件头
	heicHeader := []byte{0x66, 0x74, 0x79, 0x70, 0x68, 0x65, 0x69, 0x63}
	// AVIF 文件头
	avifHeader := []byte{0x66, 0x74, 0x79, 0x70, 0x61, 0x76, 0x69, 0x66}

	// 检查 JPEG (JPG) 文件头
	if bytes.HasPrefix(data, jpegHeader) {
		return "jpeg"
	}
	// 检查 PNG 文件头
	if bytes.HasPrefix(data, pngHeader) {
		return "png"
	}
	// 检查 GIF 文件头
	if bytes.HasPrefix(data, gifHeader1) || bytes.HasPrefix(data, gifHeader2) {
		return "gif"
	}
	// 检查 BMP 文件头
	if bytes.HasPrefix(data, bmpHeader) {
		return "bmp"
	}
	// 检查 TIFF 文件头
	if bytes.HasPrefix(data, tiffHeader1) || bytes.HasPrefix(data, tiffHeader2) {
		return "tiff"
	}
	// 检查 WebP 文件头
	if bytes.HasPrefix(data, webpHeader) {
		return "webp"
	}
	// 检查 ICO 文件头
	if bytes.HasPrefix(data, icoHeader) {
		return "ico"
	}
	// 检查 SVG 文件头
	if bytes.HasPrefix(data, svgHeader) {
		return "svg"
	}
	// 检查 HEIC 文件头
	if bytes.HasPrefix(data, heicHeader) {
		return "heic"
	}
	// 检查 AVIF 文件头
	if bytes.HasPrefix(data, avifHeader) {
		return "avif"
	}

	// 如果没有匹配到任何已知的文件头，返回空字符串
	return ""
}
