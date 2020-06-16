package modules

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

type ClusterCache struct {
	driver *cache.Cache
}

func NewClusterCache(defaultExpiration, cleanupInterval time.Duration) *ClusterCache {
	return &ClusterCache{
		driver: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *ClusterCache) AddCluster(cluster *Cluster, d time.Duration) error {
	return c.driver.Add(fmt.Sprintf("%d", cluster.ID), cluster, d)
}

func (c *ClusterCache) GetCluster(id uint64) (*Cluster, bool) {
	obj, exist := c.driver.Get(fmt.Sprintf("%d", id))
	if !exist {
		return nil, false
	}
	if cluster, ok := obj.(*Cluster); ok {
		return cluster, true
	}
	return nil, false
}
