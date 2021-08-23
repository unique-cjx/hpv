package context

import "sync"

type Container struct {
	items *sync.Map
}

// NewContainer _
func NewContainer() *Container {
	container := new(Container)
	container.items = new(sync.Map)
	return container
}

// Set _
func (c *Container) Set(key string, val interface{}) *Container {
	c.items.Store(key, val)
	return c
}

// Get _
func (c Container) Get(key string) interface{} {
	val, ok := c.items.Load(key)
	if !ok {
		return nil
	}
	return val
}
