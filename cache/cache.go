package cache

import "errors"

var ErrNoFound = errors.New("no found in cache")

type cache struct {
	m map[string]interface{}
}

func NewCache() *cache {
	return &cache{
		m: make(map[string]interface{}),
	}
}

func (c *cache) Set(key string, val interface{}) error {
	c.m[key] = val
	return nil
}

func (c *cache) Get(key string) (interface{}, error) {
	if v, ok := c.m[key]; ok {
		return v, nil
	}
	return nil, ErrNoFound
}
