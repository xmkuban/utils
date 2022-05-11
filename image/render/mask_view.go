package render

import "errors"

type MaskView struct {
	View
}

const MaskViewType = "mask"

func init() {
	RegisterViewHandlerMapping(MaskViewType, func() Viewer {
		return new(MaskView)
	})
}

func (view *MaskView) Draw(graphic *Graphic) error {
	view.drawBorder(graphic)
	graphic.SetHexColor(view.data.Color)
	floatW, floatH := float64(view.data.W), float64(view.data.H)
	floatX, floatY := float64(view.data.X), float64(view.data.Y)
	graphic.DrawRoundedRectangle(floatX, floatY, floatW, floatH, view.data.Radius)
	graphic.Fill()
	return nil
}

func (view *MaskView) Init(data *ViewData) (err error) {
	view.data = data
	return
}

func (view *MaskView) Measure() (info MeasureInfo, err error) {
	if view.data.W > 0 && view.data.H > 0 {
		actualX, actualY := view.computeBorderXY(view.data.W, view.data.H)
		view.bounds = MeasureInfo{actualX, actualY, view.data.W, view.data.H}
	} else {
		err = errors.New("invalid width or height, will ignore mask")
	}
	return view.bounds, err
}
