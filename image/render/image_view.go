package render

import (
	"bytes"
	"errors"
	"image"
	"strings"

	"github.com/go-resty/resty/v2"

	"io/ioutil"

	"github.com/xmkuban/utils/utils"
)

// Max image size: 200KB
const MAX_IMAGE_SIZE = 1024 * 512
const ImageViewType = "image"

type ImageView struct {
	View
	img image.Image
}

func init() {
	RegisterViewHandlerMapping(ImageViewType, func() Viewer {
		return new(ImageView)
	})
}

func (view *ImageView) Draw(graphic *Graphic) error {
	view.drawBorder(graphic)

	if view.img == nil {
		return errors.New("invalid image resource")
	}
	if view.data.Radius > 0 {
		bounds := view.bounds
		floatW, floatH := float64(bounds.W), float64(bounds.H)
		graphic.DrawRoundedRectangle(float64(bounds.X), float64(bounds.Y), floatW, floatH, view.data.Radius)
		graphic.Clip()
		defer graphic.ResetClip()
	}

	graphic.DrawImageAnchored(view.img, view.data.X, view.data.Y, view.data.AX, view.data.AY)
	return nil
}

func (view *ImageView) Init(data ViewData) (err error) {
	view.data = data
	if data.URL == "" {
		return errors.New("invalid image data")
	}

	if view.img == nil {
		var imageData []byte
		if len(view.data.ImageByte) > 0 {
			imageData = view.data.ImageByte
		} else {
			if strings.Contains(view.data.URL, "http") {
				var resp *resty.Response
				resp, err = resty.New().R().Get(view.data.URL)
				imageData = resp.Body()
			} else {
				imageData, err = ioutil.ReadFile(view.data.URL)
			}
		}
		if err != nil {
			return
		}
		view.img, _, err = image.Decode(bytes.NewReader(imageData))
		if err != nil {
			return err
		}
	}
	view.img = ResizeImage(view.img, uint(data.W), uint(data.H))
	return
}

func (view *ImageView) Measure() (info MeasureInfo, err error) {
	if view.img == nil {
		return info, errors.New("invalid image resource")
	}
	bounds := view.img.Bounds()
	min, max := bounds.Min, bounds.Max
	width := utils.Condition(view.data.W == 0, max.X-min.X, view.data.W).(int)
	height := utils.Condition(view.data.H == 0, max.Y-min.Y, view.data.H).(int)
	actualX, actualY := view.computeBorderXY(width, height)
	view.bounds = MeasureInfo{actualX, actualY, width, height}
	return view.bounds, nil

}
