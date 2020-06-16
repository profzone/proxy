package modules

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

type APICache struct {
	driver *cache.Cache
}

func NewAPICache(defaultExpiration, cleanupInterval time.Duration) *APICache {
	return &APICache{
		driver: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *APICache) AddAPI(api *API, d time.Duration) error {
	return c.driver.Add(fmt.Sprintf("%d", api.ID), api, d)
}

func (c *APICache) GetAPI(id uint64) (*API, bool) {
	obj, exist := c.driver.Get(fmt.Sprintf("%d", id))
	if !exist {
		return nil, false
	}
	if api, ok := obj.(*API); ok {
		return api, true
	}
	return nil, false
}
