package render

import (
	"math"

	"github.com/xmkuban/utils/utils"
)

const LineViewType = "line"

type LineView struct {
	View
}

func init() {
	RegisterViewHandlerMapping(LineViewType, func() Viewer {
		return new(LineView)
	})
}

func (view *LineView) Draw(graphic *Graphic) error {
	view.drawBorder(graphic)
	graphic.SetHexColor(view.data.Color)
	if view.data.LineW > 0 {
		graphic.SetLineWidth(view.data.LineW)
	}
	floatX, floatY := float64(view.data.X), float64(view.data.Y)
	floatX2, floatY2 := float64(view.data.X2), float64(view.data.Y2)
	graphic.DrawLine(floatX, floatY, floatX2, floatY2)
	graphic.Stroke()
	return nil
}

func (view *LineView) Init(data ViewData) (err error) {
	view.data = data
	return
}

func (view *LineView) Measure() (info MeasureInfo, err error) {
	lineW := int(math.Abs(float64(view.data.X2 - view.data.X)))
	lineH := int(math.Abs(float64(view.data.Y2 - view.data.Y)))
	actualX := utils.MinInt(view.data.X, view.data.X2)
	actualY := utils.MinInt(view.data.Y, view.data.Y2)
	view.bounds = MeasureInfo{actualX, actualY, lineW, lineH}
	return view.bounds, err
}
