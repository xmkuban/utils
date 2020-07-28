package config

import (
	"os"
	"testing"
)

type TestEnv struct {
	A string `split_words:"true"`
	B struct {
		C []int          `split_words:"true"`
		D map[string]int `split_words:"true" required:"true"`
	}
}

type TestEnvSection struct {
	C []int                  `split_words:"true"`
	D map[string]interface{} `split_words:"true" required:"true"`
}

func TestUnmarshalEnvConfig(t *testing.T) {
	os.Setenv("TEST_A", "test")
	os.Setenv("TEST_B_C", "1,2,3")
	os.Setenv("TEST_B_D", "red:1,green:2,blue:3")

	obj := TestEnv{}

	loader := &EnvConfigLoader{}
	err := loader.LoadFromResourcePath("test", &obj)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("parse env config success %v", obj)
	}

	os.Unsetenv("TEST_A")
	os.Unsetenv("TEST_B_C")
	os.Unsetenv("TEST_B_D")
}

func TestUnmarshalEnvConfigWithSection(t *testing.T) {
	os.Setenv("TEST_A", "test")
	os.Setenv("TEST_B_C", "1,2,3")
	os.Setenv("TEST_B_D", "red:1,green:2,blue:3")

	obj := TestEnvSection{}

	loader := &EnvConfigLoader{}
	err := loader.LoadFromResourcePathWithSection("test", "b", &obj)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("parse env config success %v", obj)
	}

	os.Unsetenv("TEST_A")
	os.Unsetenv("TEST_B_C")
	os.Unsetenv("TEST_B_D")
}
