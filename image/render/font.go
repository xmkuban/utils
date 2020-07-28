package render

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/xmkuban/logger"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var supportFont = make(map[string]string)
var fontCache = make(map[string]*truetype.Font)

func FontInit(_supportFont map[string]string) {
	supportFont = _supportFont
	for k, v := range supportFont {
		f, err := loadFont(v)
		if err != nil {
			logger.Error(err)
			continue
		}
		fontCache[k] = f
	}
}

func IsExistFont(supportFontKey string, supportFontPath string) bool {
	if _, ok := supportFont[supportFontKey]; ok {
		if supportFontPath == "" {
			return true
		}
		if supportFont[supportFontKey] == supportFontPath {
			return true
		}
	}
	return false
}

func AddOrUpdateFont(supportFontKey string, supportFontPath string) error {
	if supportFontPath == "" {
		return errors.New("font not find")
	}
	if _, ok := supportFont[supportFontKey]; ok {
		if supportFont[supportFontKey] == supportFontPath {
			return nil
		}
	}
	supportFont[supportFontKey] = supportFontPath
	f, err := loadFont(supportFontPath)
	if err != nil {
		logger.Error(err)
		return err
	}
	fontCache[supportFontKey] = f
	return nil
}

func loadFont(path string) (f *truetype.Font, err error) {
	var fontBytes []byte
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		var resp *resty.Response
		resp, err = resty.New().R().Get(path)
		fontBytes = resp.Body()
	} else {
		fontBytes, err = ioutil.ReadFile(path)
	}
	if err != nil {
		return nil, err
	}
	f, err = freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	return
}

func loadFontFace(supportFont string, points float64) (font.Face, error) {
	if _, ok := fontCache[supportFont]; !ok {
		return nil, errors.New("font not find")
	}
	f := fontCache[supportFont]

	if f == nil {
		return nil, errors.New("font not find")
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		//Hinting: font.HintingFull,
	})
	return face, nil
}
