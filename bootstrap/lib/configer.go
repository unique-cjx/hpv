package lib

type YamlConfig struct {
	Logger  LoggerConfig  `yaml:"logger,omitempty"`
	YueMiao YueMiaoConfig `yaml:"yuemiao"`
}

type LoggerConfig struct {
	Level string `yaml:"level,omitempty" default:"debug"`
	Debug bool   `yml:"debug,omitempty" default:"true"`
}

type YueMiaoConfig struct {
	Tk       string `yaml:"tk"`
	Province Region `yaml:"province"`
	City     Region `yaml:"city"`
}

type Region struct {
	Name string `yaml:"name"`
	Code string `yaml:"code"`
}
