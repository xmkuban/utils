package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

const (
	DEFAULT_YAML_FILE_SUFFIX = ".yaml"
)

type YamlConfConfigLoader struct {
	Path string
}

func (loader *YamlConfConfigLoader) LoadFromResourcePath(path string, spec interface{}) error {
	if filepath.Ext(path) == "" {
		path = path + DEFAULT_YAML_FILE_SUFFIX
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, spec)
	return err
}

func (loader *YamlConfConfigLoader) LoadFromResourcePathWithSection(path string, section string, spec interface{}) error {
	return errors.New("not found section")
}

func (loader *YamlConfConfigLoader) LoadFromFlagPath(spec interface{}) error {
	if filepath.Ext(loader.Path) == "" {
		loader.Path = loader.Path + DEFAULT_YAML_FILE_SUFFIX
	}
	file, err := os.Open(loader.Path)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, spec)
	return err
}
