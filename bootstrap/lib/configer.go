package lib

type YamlConfig struct {
	Logger  LoggerCfg     `yaml:"logger,omitempty"`
	YM      YMCfg         `yaml:"yuemiao"`
	Regions []ProvinceCfg `yaml:"region"`
}

type LoggerCfg struct {
	Level string `yaml:"level,omitempty" default:"debug"`
	// Debug bool   `yml:"debug,omitempty" default:"true"`
}

type YMCfg struct {
	Tk string `yaml:"tk"`
}

type ProvinceCfg struct {
	Name string  `yaml:"name"`
	Code int     `yaml:"code"`
	City CityCfg `ymal:"city"`
}

type CityCfg struct {
	Name string `yaml:"name"`
	Code int    `yaml:"code"`
}
