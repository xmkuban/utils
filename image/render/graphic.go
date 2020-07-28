package render

import (
	"errors"
	"image"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
)

type Graphic struct {
	impl *gg.Context
}

func NewGraphic(context *gg.Context) (*Graphic, error) {
	if context == nil {
		return nil, errors.New("graphic implementation can not be nil")
	}
	return &Graphic{context}, nil
}

func GetImageSize(im image.Image) (w, h uint) {
	bounds := im.Bounds()
	min, max := bounds.Min, bounds.Max
	return uint(max.X - min.X), uint(max.Y - min.Y)
}

func ResizeImage(im image.Image, w, h uint) image.Image {
	if imageW, imageH := GetImageSize(im); (imageW != w || imageH != h) && (w != 0 && h != 0) {
		im = resize.Resize(w, h, im, resize.NearestNeighbor)
	}
	return im
}

func MeasureString(s string, f font.Face) (w, h float64) {
	ctx := gg.NewContext(0, 0)
	ctx.SetFontFace(f)
	return ctx.MeasureString(s)
}

func MeasureVerticalString(s string, f font.Face) (w, h float64) {
	ctx := gg.NewContext(0, 0)
	ctx.SetFontFace(f)
	return ctx.MeasureString(s)
}

func LoadImage(path string) (img image.Image, err error) {
	return gg.LoadImage(path)
}

func (g *Graphic) SetRGB(cr, cg, cb float64) {
	g.impl.SetRGB(cr, cg, cb)
}

func (g *Graphic) Push() {
	g.impl.Push()
}

func (g *Graphic) Pop() {
	g.impl.Pop()
}

func (g *Graphic) SetRGBA(cr, cg, cb, ca float64) {
	g.impl.SetRGBA(cr, cg, cb, ca)
}

func (g *Graphic) SetHexColor(hexColor string) {
	g.impl.SetHexColor(hexColor)
}

func (g *Graphic) SetLineWidth(lineWidth float64) {
	g.impl.SetLineWidth(lineWidth)
}

func (g *Graphic) Fill() {
	g.impl.Fill()
}

func (g *Graphic) FillPreserve() {
	g.impl.FillPreserve()
}

func (g *Graphic) Clip() {
	g.impl.Clip()
}

func (g *Graphic) Clear() {
	g.impl.Clear()
}

func (g *Graphic) ResetClip() {
	g.impl.ResetClip()
}

func (g *Graphic) DrawPoint(x, y, r float64) {
	g.impl.DrawPoint(x, y, r)
}

func (g *Graphic) DrawLine(x1, y1, x2, y2 float64) {
	g.impl.DrawLine(x1, y1, x2, y2)
}

func (g *Graphic) DrawRectangle(x, y, w, h float64) {
	g.impl.DrawRectangle(x, y, w, h)
}

func (g *Graphic) DrawRoundedRectangle(x, y, w, h, r float64) {
	g.impl.DrawRoundedRectangle(x, y, w, h, r)
}

func (g *Graphic) DrawEllipticalArc(x, y, rx, ry, angle1, angle2 float64) {
	g.impl.DrawEllipticalArc(x, y, rx, ry, angle1, angle2)
}

func (g *Graphic) DrawEllipse(x, y, rx, ry float64) {
	g.impl.DrawEllipse(x, y, rx, ry)
}

func (g *Graphic) DrawArc(x, y, r, angle1, angle2 float64) {
	g.impl.DrawArc(x, y, r, angle1, angle2)
}

func (g *Graphic) DrawCircle(x, y, r float64) {
	g.impl.DrawCircle(x, y, r)
}

func (g *Graphic) DrawRegularPolygon(n int, x, y, r, rotation float64) {
	g.impl.DrawRegularPolygon(n, x, y, r, rotation)
}

func (g *Graphic) DrawImage(im image.Image, x, y int) {
	g.impl.DrawImage(im, x, y)
}

func (g *Graphic) DrawImageAnchored(im image.Image, x, y int, ax, ay float64) {
	g.impl.DrawImageAnchored(im, x, y, ax, ay)
}

func (g *Graphic) SetFontFace(fontFace font.Face) {
	g.impl.SetFontFace(fontFace)
}

func (g *Graphic) LoadFontFace(path string, points float64) error {
	return g.impl.LoadFontFace(path, points)
}

func (g *Graphic) DrawString(s string, x, y float64) {
	g.impl.DrawString(s, x, y)
}

func (g *Graphic) DrawStringAnchored(s string, x, y, ax, ay float64) {
	g.impl.DrawStringAnchored(s, x, y, ax, ay)
}

func (g *Graphic) DrawStringVertical(s string, x, y, ax, ay float64) {
	g.impl.DrawStringAnchored(s, x, y, ax, ay)
}

func (g *Graphic) DrawStringWrapped(s string, x, y, ax, ay, width, lineSpacing float64, align int) {
	g.impl.DrawStringWrapped(s, x, y, ax, ay, width, lineSpacing, gg.Align(align))
}

func (g *Graphic) Stroke() {
	g.impl.Stroke()
}

func (g *Graphic) WordWrap(s string, w float64) []string {
	return g.impl.WordWrap(s, w)
}

func (g *Graphic) Translate(x, y float64) {
	g.impl.Translate(x, y)
}

func (g *Graphic) Scale(x, y float64) {
	g.impl.Scale(x, y)
}

func (g *Graphic) ScaleAbout(sx, sy, x, y float64) {
	g.impl.ScaleAbout(sx, sy, x, y)
}

func (g *Graphic) Rotate(angle float64) {
	g.impl.Rotate(angle)
}

func (g *Graphic) RotateAbout(angle, x, y float64) {
	g.impl.RotateAbout(angle, x, y)
}
