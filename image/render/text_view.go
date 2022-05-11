package render

import (
	"errors"

	"github.com/xmkuban/utils/utils"
	"golang.org/x/image/font"
)

type Orientation int

type TextView struct {
	View
	fontFace   font.Face
	isVertical bool //是否是竖排文字
}

const (
	Horizontal Orientation = iota //横排
	Vertical                      //竖排
)

const TextViewType = "text"

func init() {
	RegisterViewHandlerMapping(TextViewType, func() Viewer {
		return new(TextView)
	})
}

func (view *TextView) Draw(graphic *Graphic) error {
	view.drawBorder(graphic)
	viewData := view.data
	x, y := float64(viewData.X), float64(viewData.Y)
	graphic.SetHexColor(viewData.Color)
	graphic.SetFontFace(view.fontFace)
	if view.isVertical {
		graphic.DrawStringVertical(viewData.Text, x, y, viewData.AY, viewData.AY)
	} else {
		graphic.DrawStringAnchored(viewData.Text, x, y, viewData.AX, viewData.AY)
	}
	return nil
}

func (view *TextView) adjustFontSize(maxSize int, face font.Face, measure func(text string, face font.Face) (size float64)) (f font.Face, err error) {
	viewData := view.data
	floatSize, fontSize := float64(maxSize), viewData.Size
	f = face
	var measureSize float64
	for {
		measureSize = measure(viewData.Text, f)
		if measureSize <= floatSize {
			break
		} else {
			fontSize -= 2
			f, err = loadFontFace(viewData.Font, fontSize)
			if err != nil || fontSize <= 2 {
				err = errors.New("invalid text size")
				return
			}
		}
	}
	return
}

func (view *TextView) Init(data *ViewData) (err error) {
	view.data = data
	if data.Text == "" {
		err = errors.New("empty text for text view")
		return
	}
	view.data.Text = utils.FilterEmoji(view.data.Text)
	view.isVertical = data.Orientation == Vertical
	fontFace, err := loadFontFace(data.Font, data.Size)
	if view.isVertical && data.H > 0 {
		fontFace, err = view.adjustFontSize(data.H, fontFace, func(text string, face font.Face) (size float64) {
			_, h := MeasureVerticalString(text, face)
			return h
		})
	} else if data.W > 0 {
		fontFace, err = view.adjustFontSize(data.W, fontFace, func(text string, face font.Face) (size float64) {
			w, _ := MeasureString(text, face)
			return w
		})
	}
	if err == nil {
		view.fontFace = fontFace
	}
	return
}

func (view *TextView) Measure() (info MeasureInfo, err error) {
	viewData := view.data
	var w, h float64
	if view.isVertical {
		w, h = MeasureVerticalString(viewData.Text, view.fontFace)
	} else {
		w, h = MeasureString(viewData.Text, view.fontFace)
	}
	actualX, actualY := view.computeBorderXY(int(w), int(h))
	if !view.isVertical {
		actualY -= int(h - 5)
	}
	view.bounds = MeasureInfo{actualX, actualY, int(w), int(h)}
	return view.bounds, nil
}
