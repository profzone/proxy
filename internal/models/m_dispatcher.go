package models

import (
	"bytes"
	"encoding/gob"
	"github.com/profzone/eden-framework/pkg/timelib"
)

type BreakerConf struct {
	// Half-open状态下最多能进入的请求数量
	MaxRequests uint32 `json:"maxRequests"`
	// Close状态下重置内部统计的时间
	Interval timelib.DurationString `json:"interval"`
	// Open状态下变更为Half-open状态的时间
	Timeout timelib.DurationString `json:"timeout"`
}

type Dispatcher struct {
	// 路由
	Router *Router `json:"router,omitempty" default:""`
	// 熔断器
	BreakerConf *BreakerConf `json:"breaker,omitempty" default:""`

	WriteTimeout timelib.DurationString `json:"writeTimeout" default:""`
	ReadTimeout  timelib.DurationString `json:"readTimeout" default:""`
	ClusterID    uint64                 `json:"clusterID,string"`
}

func (d *Dispatcher) GobDecode(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(d)
	return err
}
