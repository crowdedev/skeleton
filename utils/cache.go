package utils

import (
	"fmt"
	"time"

	configs "github.com/crowdeco/skeleton/configs"
	cachita "github.com/gadelkareem/cachita"
)

type cache struct {
	pool cachita.Cache
}

func NewCache() *cache {
	return &cache{
		pool: cachita.Memory(),
	}
}

func (c *cache) Set(key string, value interface{}) {
	err := c.pool.Put(key, value, time.Duration(configs.Env.CacheLifetime)*time.Second)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *cache) Get(key string) (interface{}, bool) {
	var data interface{}
	err := c.pool.Get(key, &data)
	if err != nil {
		return nil, false
	}

	return &data, true
}
