package db

import "longhorn/proxy/internal/modules"

type DB interface {
	CreateCluster(c modules.Cluster)
}
