package utils

import (
	"fmt"
	"time"

	configs "github.com/crowdeco/skeleton/configs"
	cachita "github.com/gadelkareem/cachita"
)

type Cache struct {
	Env  *configs.Env
	Pool cachita.Cache
}

func (c *Cache) Set(key string, value interface{}) {
	err := c.Pool.Put(key, &value, time.Duration(c.Env.CacheLifetime)*time.Second)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Cache) Invalidate(key string) {
	err := c.Pool.Invalidate(key)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	var data interface{}
	err := c.Pool.Get(key, &data)
	if err != nil {
		return nil, false
	}

	return data, true
}
