package render

import (
	qrcode "github.com/skip2/go-qrcode"
)

type QrView struct {
	ImageView
}

const QrViewType = "qr"

func init() {
	RegisterViewHandlerMapping(QrViewType, func() Viewer {
		return new(QrView)
	})
}

func (view *QrView) Init(data *ViewData) error {
	view.data = data
	qr, err := qrcode.New(data.Text, qrcode.RecoveryLevel(data.Level))
	if err != nil {
		return err
	}
	view.img = qr.Image(int(data.Size))
	return nil
}
