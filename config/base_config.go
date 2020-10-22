package config

import (
	"github.com/xmkuban/logger"

	"context"
	"os"
)

type ServiceCallBaseConfig struct {
	ForceHttps bool
}

type StorageBaseConfig struct {
	Type string `yaml:"type" split_words:"true"`
	//credential key - credential map
	Config map[string]interface{} `yaml:"config" split_words:"true"`
}

type BaseConfig struct {
	Debug             bool
	ConfigType        string
	ConfigPath        string
	AppName           string
	AppVersion        string
	AppVersionCode    int64
	AppBuildTimestamp int64
	Env               string
	ctx               context.Context

	ServiceCfg *ServiceCallBaseConfig

	//use_for- countryShort - config map
	StorageCfg map[string]map[string]*StorageBaseConfig
}

func (c *BaseConfig) SetCTX(ctx context.Context) {
	c.ctx = ctx
}

func (c *BaseConfig) Extra(key string, val interface{}) {
	if c.ctx == nil {
		c.ctx = context.Background()
	}

	c.ctx = context.WithValue(c.ctx, key, val)
}

func (c *BaseConfig) GetExtra(key string) (val interface{}, hit bool) {
	if c.ctx == nil {
		c.ctx = context.Background()
		return nil, false
	}

	val = c.ctx.Value(key)
	if val == nil {
		return val, false
	}
	return val, true
}

func (c *BaseConfig) GetExtraString(key string) (val string, hit bool) {
	var valRaw interface{}
	valRaw, hit = c.GetExtra(key)
	if hit {
		val, hit = valRaw.(string)
		return
	} else {
		return "", false
	}
}

func (c *BaseConfig) LoadCommon() error {
	loader := &YamlConfConfigLoader{
		Path: "./config.yaml",
	}
	return loader.LoadFromResourcePathWithSection(c.ConfigPath, "common.storage", &c.StorageCfg)
}

var PresetConfig *BaseConfig

type ConfigLoader interface {
	LoadFromResourcePath(path string, spec interface{}) error
	LoadFromResourcePathWithSection(path string, section string, spec interface{}) error
	LoadFromFlagPath(spec interface{}) error
}

func GetConfigLoaderByType(loaderType string) ConfigLoader {
	if loaderType == "env" {
		return &EnvConfigLoader{}
	} else if loaderType == "yaml" {
		return &YamlConfConfigLoader{}
	}
	return &GoConfConfigLoader{}
}

func GetConfigLoader() ConfigLoader {
	if PresetConfig == nil {
		logger.Error("Preset config need to be set")
		panic("preset config is nil")
		return nil
	}
	if PresetConfig.ConfigType == "env" {
		return &EnvConfigLoader{
			Path: PresetConfig.ConfigPath,
		}
	} else if PresetConfig.ConfigType == "yaml" {
		return &YamlConfConfigLoader{
			Path: PresetConfig.ConfigPath,
		}
	} else {
		return &GoConfConfigLoader{
			Path: PresetConfig.ConfigPath,
		}
	}
}

func InitPresetConfig() {
	PresetConfig = new(BaseConfig)
	PresetConfig.ConfigType = "goconf"
	PresetConfig.ConfigPath = "config"
	PresetConfig.ServiceCfg = new(ServiceCallBaseConfig)

	//force use https
	httpsEnableStr, ok := os.LookupEnv("DISABLE_FORCE_HTTPS")
	if !ok || httpsEnableStr == "" || httpsEnableStr == "0" {
		PresetConfig.ServiceCfg.ForceHttps = true
	} else {
		PresetConfig.ServiceCfg.ForceHttps = false
	}
	PresetConfig.StorageCfg = make(map[string]map[string]*StorageBaseConfig)
}

func InitPresetConfigWithPath(path string) {
	PresetConfig = new(BaseConfig)
	PresetConfig.ConfigType = "goconf"
	PresetConfig.ConfigPath = path
	PresetConfig.ServiceCfg = new(ServiceCallBaseConfig)

	//force use https
	httpsEnableStr, ok := os.LookupEnv("DISABLE_FORCE_HTTPS")
	if !ok || httpsEnableStr == "" || httpsEnableStr == "0" {
		PresetConfig.ServiceCfg.ForceHttps = true
	} else {
		PresetConfig.ServiceCfg.ForceHttps = false
	}

	//init storage config
	PresetConfig.StorageCfg = make(map[string]map[string]*StorageBaseConfig)
}

func EasyLoadConfig(spec interface{}) error {
	loader := GetConfigLoader()

	return loader.LoadFromFlagPath(spec)
}
