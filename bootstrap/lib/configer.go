package lib

type YamlConfig struct {
	Logger LoggerCfg `yaml:"logger,omitempty"`
	YM     YMCfg     `yaml:"yuemiao"`
	Region []Region
}

type LoggerCfg struct {
	Level string `yaml:"level,omitempty" default:"debug"`
	Debug bool   `yml:"debug,omitempty" default:"true"`
}

type YMCfg struct {
	Tk string `yaml:"tk"`
}

type Region struct {
	Name string `yaml:"name"`
	Code int    `yaml:"code"`
}
