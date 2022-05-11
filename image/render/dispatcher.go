package render

import (
	"bytes"
	"errors"
	"image/jpeg"
	"sync"

	"github.com/xmkuban/logger"

	"github.com/fogleman/gg"
	"github.com/xmkuban/utils/utils"
)

type ImageFormat int

const (
	JPEG ImageFormat = iota
	PNG
)

var viewHandlerMappingStore = make(map[string]func() Viewer, 0)

func RegisterViewHandlerMapping(key string, viewMapping func() Viewer) {
	if key != "" && viewMapping != nil {
		viewHandlerMappingStore[key] = viewMapping
	}
}

func getBounds(viewer Viewer) (width, height int, err error) {
	if viewer == nil {
		err = errors.New("invalid viewer")
		return
	}
	bounds, err := viewer.Measure()
	if err != nil {
		return
	}
	width = bounds.X + bounds.W
	height = bounds.Y + bounds.H
	return
}

func Dispatch(renderRequest ImageData) (renderResult []byte, err error) {
	var actualWidth, actualHeight int
	views := renderRequest.View
	viewOrder := make([]Viewer, len(views))
	defer func() {
		if viewOrder != nil {
			viewOrder = nil
		}
	}()
	waitInit := sync.WaitGroup{}
	waitInit.Add(len(views))
	for i, value := range views {
		go (func(idx int, viewData ViewData) {
			defer waitInit.Done()
			viewMapping := viewHandlerMappingStore[viewData.Type]
			if viewMapping == nil {
				return
			}
			viewer := viewMapping()
			err = viewer.Init(viewData)
			if err != nil {
				return
			}
			viewW, viewH, err := getBounds(viewer)
			if err != nil {
				return
			}
			actualWidth, actualHeight = utils.MaxInt(actualWidth, viewW), utils.MaxInt(actualHeight, viewH)
			viewOrder[idx] = viewer
		})(i, value)
	}
	waitInit.Wait()

	if renderRequest.W > 0 {
		actualWidth = renderRequest.W
	}
	if renderRequest.H > 0 {
		actualHeight = renderRequest.H
	}
	ctx := gg.NewContext(actualWidth, actualHeight)
	if renderRequest.Color != "" {
		ctx.SetHexColor(renderRequest.Color)
		ctx.Clear()
	}
	graphic, err := NewGraphic(ctx)
	defer func() {
		if graphic != nil {
			graphic = nil
		}
		if ctx != nil {
			ctx.ClearPath()
			ctx = nil
		}
	}()
	if err != nil {
		return
	}
	for _, viewer := range viewOrder {
		if viewer != nil {
			err = viewer.Draw(graphic)
			if err != nil {
				logger.Error(err)
			}
		}
	}
	buf := new(bytes.Buffer)
	switch renderRequest.Format {
	case PNG:
		err = ctx.EncodePNG(buf)
	default:
		quality := utils.Condition(renderRequest.Quality <= 0, 75, renderRequest.Quality).(int)
		err = ctx.EncodeJPG(buf, &jpeg.Options{Quality: quality})
	}

	if err != nil {
		return
	}
	return buf.Bytes(), nil
}
