package render

import "image"

type ImageData struct {
	Quality int         `json:"quality"`
	Format  ImageFormat `json:"format"`
	View    []ViewData  `json:"views"`
}

type ViewData struct {
	// Common
	W           int     `json:"w"`
	H           int     `json:"h"`
	X           int     `json:"x"`
	Y           int     `json:"y"`
	AX          float64 `json:"ax"`
	AY          float64 `json:"ay"`
	BorderWidth float64 `json:"borderWidth"`
	BorderColor string  `json:"borderColor"`
	Padding     float64 `json:"padding"`
	Type        string  `json:"type"`

	// ImageView
	URL       string  `json:"url"`
	Radius    float64 `json:"radius"`
	ImageByte []byte  `json:"image_byte"`

	// TextView
	Text        string      `json:"text"`
	Font        string      `json:"font"`
	FontPath    string      `json:"font_path"`
	Size        float64     `json:"size"` // qr code 也会用到
	Color       string      `json:"color"`
	Orientation Orientation `json:"orientation"`

	// QRView
	Level int `json:"level"`

	// LineView
	LineW float64 `json:"lineW"`
	X2    int     `json:"x2"`
	Y2    int     `json:"y2"`

	IsCache bool `json:"is_cache"`

	img image.Image `json:"-"`
}

func (this *ViewData) SetImg(img image.Image) {
	this.img = img
}
