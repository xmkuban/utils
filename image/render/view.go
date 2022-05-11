package render

type Viewer interface {
	Init(data *ViewData) (err error)
	Measure() (info MeasureInfo, err error)
	Draw(graphic *Graphic) (err error)
}

type MeasureInfo struct {
	X, Y, W, H int
}

type View struct {
	data   *ViewData
	bounds MeasureInfo
}

func (v *View) computeBorderXY(w, h int) (actualX, actualY int) {
	return v.data.X - int(float64(w)*v.data.AX), v.data.Y - int(float64(h)*v.data.AY)
}

func (v *View) drawBorder(graphic *Graphic) {
	if v.data.BorderWidth > 0 {
		borderWidth := v.data.BorderWidth
		borderColor := v.data.BorderColor
		innerRect := v.bounds
		padding := v.data.Padding
		innerX, innerY := float64(innerRect.X), float64(innerRect.Y)
		innerW, innerH := float64(innerRect.W), float64(innerRect.H)

		graphic.SetRGBA(0, 0, 0, 0)
		graphic.DrawRoundedRectangle(innerX-padding, innerY-padding, innerW+padding*2, innerH+padding*2, v.data.Radius)
		graphic.SetLineWidth(borderWidth)
		graphic.FillPreserve()

		if borderColor == "" {
			graphic.SetRGB(0, 0, 0)
		} else {
			graphic.SetHexColor(borderColor)
		}

		graphic.Stroke()
	}
}
