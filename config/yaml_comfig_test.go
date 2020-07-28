package config

import "testing"

type TestYaml struct {
	A string `yaml:"a"`
	B struct {
		C []int                  `yaml:"c"`
		D map[string]interface{} `yaml:"d"`
	} `yaml:"b"`
}

type TestYamlSection struct {
	C []int                  `yaml:"c"`
	D map[string]interface{} `yaml:"d"`
}

func TestUnmarshalYamlConfig(t *testing.T) {
	obj := TestYaml{}

	loader := &YamlConfConfigLoader{}
	err := loader.LoadFromResourcePath("./test.yaml", &obj)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("parse yaml config success %v", obj)
	}
}

func TestUnmarshalYamlConfigWithSection(t *testing.T) {
	obj := TestYamlSection{}

	loader := &YamlConfConfigLoader{}
	err := loader.LoadFromResourcePathWithSection("./test.yaml", "b", &obj)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("parse yaml config success %v", obj)
	}
}
