package config

import (
	"errors"
	"path/filepath"

	"github.com/Terry-Mao/goconf"
)

const (
	DEFAULT_FILE_SUFFIX = ".conf"
)

type GoConfConfigLoader struct {
	Path string
}

func (loader *GoConfConfigLoader) LoadFromResourcePath(path string, spec interface{}) error {
	if filepath.Ext(path) == "" {
		path = path + DEFAULT_FILE_SUFFIX
	}
	conf := goconf.New()
	err := conf.Parse(path)
	if err != nil {
		return err
	}
	err = conf.Unmarshal(spec)
	//err := goconf.NewConf(path, spec)
	return err
}

func (loader *GoConfConfigLoader) LoadFromResourcePathWithSection(path string, section string, spec interface{}) error {
	return errors.New("unsupported function")
}

func (loader *GoConfConfigLoader) LoadFromFlagPath(spec interface{}) error {
	conf := goconf.New()
	if filepath.Ext(loader.Path) == "" {
		loader.Path = loader.Path + DEFAULT_FILE_SUFFIX
	}
	err := conf.Parse(loader.Path)
	if err != nil {
		return err
	}
	err = conf.Unmarshal(spec)
	//err := goconf.NewConf(path, spec)
	return err
}
