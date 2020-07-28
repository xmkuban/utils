package config

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type EnvConfigLoader struct {
	Path string
}

func (loader *EnvConfigLoader) LoadFromResourcePath(path string, spec interface{}) error {
	err := envconfig.Process(path, spec)
	return err
}

func (loader *EnvConfigLoader) LoadFromResourcePathWithSection(path string, section string, spec interface{}) error {
	var prefix string
	sectionArr := strings.Split(section, ".")
	section = strings.Join(sectionArr, "_")
	if path == "" {
		prefix = section
	} else {
		prefix = fmt.Sprintf("%s_%s", path, section)
	}
	err := envconfig.Process(prefix, spec)
	return err
}

func (loader *EnvConfigLoader) LoadFromFlagPath(spec interface{}) error {
	return loader.LoadFromResourcePath(loader.Path, spec)
}
