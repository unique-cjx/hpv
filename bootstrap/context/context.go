package context

import "hpv/bootstrap/lib"

type Context struct {
	*Container
}

// NewContext _
func NewContext() *Context {
	context := new(Context)
	context.Container = NewContainer()
	return context
}

// GetAppConfig _
func (c Context) GetAppConfig() (conf *lib.YamlConfig) {
	tmpConf := c.Get("app_config")
	conf, ok := tmpConf.(*lib.YamlConfig)
	if !ok {
		conf = nil
	}
	return
}

// SetAppConfig _
func (c Context) SetAppConfig(conf *lib.YamlConfig) {
	c.Set("app_config", conf)
}
