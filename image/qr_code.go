package image

import (
	"image"

	"io"

	"image/jpeg"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/xmkuban/logger"
)

func CreateQrCode(content string, size int) (img image.Image, err error) {
	qr, err := qrcode.New(content, qrcode.Highest)
	if err != nil {
		logger.Error(err)
		return
	}
	img = qr.Image(size)
	return
}

func CreateQrCodeAndJPEGWrite(content string, size int, quality int, w io.Writer) error {

	img, err := CreateQrCode(content, size)
	if err != nil {
		return err
	}

	err = jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	return err
}
