package modules

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var ClusterContainer *ClusterCache
var APIContainer *APICache

func init() {
	APIContainer = NewAPICache(cache.NoExpiration, cache.NoExpiration)
	ClusterContainer = NewClusterCache(5*time.Minute, 10*time.Minute)
}
