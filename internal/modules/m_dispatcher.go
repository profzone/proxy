package modules

import "time"

type Dispatcher struct {
	Router Router `json:"router" default:""`
	// Validations
	WriteTimeout time.Duration `json:"writeTimeout" default:""`
	ReadTimeout  time.Duration `json:"readTimeout" default:""`
	ClusterID    uint64        `json:"clusterID,string"`
}
