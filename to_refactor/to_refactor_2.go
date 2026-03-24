package to_refactor

import "sync"

/*
explain and refactor if needed :D

type cache struct {
	mutex sync.Mutex
	data  map[string]string
}

func newCache() *cache {
	return &cache{
		data: make(map[string]string),
	}
}

func (c *cache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = value
}

func (c *cache) Get(key string) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.Size() > 0 {
		return c.data[key]
	}

	return ""
}

func (c *cache) Size() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return len(c.data)
}
*/

type cache struct {
	mutex sync.Mutex
	data  map[string]string
}

func newCache() *cache {
	return &cache{
		data: make(map[string]string),
	}
}

func (c *cache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = value
}

func (c *cache) Get(key string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.sizeLocked() > 0 {
		v, ok := c.data[key]
		return v, ok
	}
	return "", false
}

func (c *cache) Size() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.sizeLocked()
}

func (c *cache) sizeLocked() int {
	return len(c.data)
}
