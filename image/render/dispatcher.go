package render

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/fogleman/gg"
	"github.com/go-resty/resty/v2"
	"github.com/xmkuban/logger"
	"github.com/xmkuban/utils/cache"
	"github.com/xmkuban/utils/utils"
)

var imageCache *cache.MemoryCache

func init() {
	imageCache = cache.NewMemoryCache(nil)
}

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

func Dispatch(renderRequest *ImageData) (renderResult []byte, err error) {
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
	actualLock := sync.Mutex{}
	for _i, _value := range views {
		value := _value
		i := _i
		go func() {
			defer waitInit.Done()
			viewer, viewW, viewH, err1 := initView(value)
			if err1 != nil {
				err = err1
				return
			}
			if viewW > 0 || viewH > 0 {
				actualLock.Lock()
				actualWidth = utils.MaxInt(actualWidth, viewW)
				actualHeight = utils.MaxInt(actualHeight, viewH)
				actualLock.Unlock()
			}
			viewOrder[i] = viewer
		}()
	}
	waitInit.Wait()
	if err != nil {
		return nil, err
	}

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
		if viewer == nil {
			continue
		}
		err = viewer.Draw(graphic)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}
	buf := new(bytes.Buffer)
	switch renderRequest.Format {
	case PNG:
		err = ctx.EncodePNG(buf)
	default:
		quality := utils.Condition(renderRequest.Quality <= 0, 90, renderRequest.Quality).(int)
		err = ctx.EncodeJPG(buf, &jpeg.Options{Quality: quality})
	}

	if err != nil {
		return
	}
	return buf.Bytes(), nil
}

func initView(viewData *ViewData) (Viewer, int, int, error) {
	viewMapping := viewHandlerMappingStore[viewData.Type]
	if viewMapping == nil {
		return nil, 0, 0, nil
	}

	if viewData.Type == TextViewType && viewData.Font != "" {
		if !IsExistFont(viewData.Font, viewData.FontPath) {
			err1 := AddOrUpdateFont(viewData.Font, viewData.FontPath)
			if err1 != nil {
				return nil, 0, 0, err1
			}
		}
	}

	if viewData.Type == ImageViewType {
		var err1 error
		viewData, err1 = imageView(viewData)
		if err1 != nil {
			return nil, 0, 0, err1
		}
	}

	viewer := viewMapping()
	err1 := viewer.Init(viewData)
	if err1 != nil {
		return nil, 0, 0, err1
	}
	viewW, viewH, err1 := getBounds(viewer)
	if err1 != nil {
		return nil, 0, 0, err1
	}
	return viewer, viewW, viewH, nil
}

func imageView(v *ViewData) (*ViewData, error) {
	cacheKey := ""
	isCache := v.IsCache
	if v.URL != "" {
		if isCache {
			cacheKey = "image_" + utils.MD5(v.URL)
		}
	} else {
		isCache = false
	}
	if isCache {
		cacheData := imageCache.Get(cacheKey)
		if cacheData != nil {
			return cacheData.(*ViewData), nil
		}
	}
	if v.GetImg() != nil {
		return v, nil
	}

	var err error
	if len(v.ImageByte) == 0 {
		if v.URL == "" {
			return nil, errors.New("param error")
		} else if strings.Contains(v.URL, "http://") || strings.Contains(v.URL, "https://") {
			resp, err1 := resty.New().R().Get(v.URL)
			if err1 != nil {
				err = err1
				logger.Error(err1)
				return nil, err
			}
			v.ImageByte = resp.Body()
		} else {
			v.ImageByte, err = ioutil.ReadFile(v.URL)
			if err != nil {
				logger.Errorf("url:%s,err:%s", v.URL, err)
				return nil, err
			}
		}
	}

	var img image.Image
	img, _, err = image.Decode(bytes.NewReader(v.ImageByte))

	img = ResizeImage(img, uint(v.W), uint(v.H))

	v.SetImg(img)

	if isCache {
		imageCache.Put(cacheKey, v, time.Hour*4)
	}
	return v, nil
}
