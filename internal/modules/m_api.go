package modules

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/juju/ratelimit"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"longhorn/proxy/internal/constants/enum"
	"longhorn/proxy/internal/global"
	"longhorn/proxy/internal/storage"
	"time"
)

type BreakerConf struct {
	// Half-open状态下最多能进入的请求数量
	MaxRequests uint32
	// Close状态下重置内部统计的时间
	Interval time.Duration
	// Open状态下变更为Half-open状态的时间
	Timeout time.Duration
}

type APIContainer map[uint64]*API

func (c APIContainer) Add(api *API) {
	c[api.ID] = api
}

func (c APIContainer) Get(id uint64) *API {
	return c[id]
}

type API struct {
	// 唯一标识
	ID uint64 `json:"id" default:""`
	// 接口名称
	Name string `json:"name"`
	// 接口URL匹配模式
	URLPattern string `json:"urlPattern"`
	// 接口匹配方法
	Method string `json:"method" default:""`
	// 接口状态
	Status enum.ApiStatus `json:"status" default:"UP"`
	// IP黑白名单 format: <blacklist(>ip[,]...<)whitelist(>ip[,]...<)>
	IPControl    string `json:"ipControl,omitempty" default:""`
	ipController *IPController
	// 最大QPS
	MaxQPS  int64 `json:"maxQPS,omitempty" default:""`
	limiter *ratelimit.Bucket
	// TODO Validations
	// 反向代理调度
	Dispatchers []Dispatcher `json:"dispatcher"`
	// TODO Fusion
}

func (v *API) SetIdentity(id uint64) {
	v.ID = id
}

func (v API) GetIdentity() uint64 {
	return v.ID
}

func (v API) Marshal() (result []byte, err error) {
	buf := bytes.NewBuffer(result)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(v)
	return buf.Bytes(), err
}

func (v *API) Unmarshal(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(v)
	if err != nil {
		return
	}

	if v.IPControl != "" {
		v.ipController, err = newIPController(v.IPControl)
		if err != nil {
			return
		}
	}
	if v.MaxQPS > 0 {
		v.limiter = ratelimit.NewBucket(time.Second/time.Duration(v.MaxQPS), v.MaxQPS)
	}
	return
}

func (v *API) WalkDispatcher(walking func(dispatcher *Dispatcher) error) {
	for _, d := range v.Dispatchers {
		err := walking(&d)
		if err != nil {
			// TODO change error display
			logrus.Error(err)
		}
	}
}

func (v *API) FilterIPControl(req *fasthttp.Request) bool {
	if v.ipController != nil {
		return v.ipController.filter(req)
	}
	return true
}

func (v *API) FilterQPS() bool {
	if v.limiter != nil && v.limiter.TakeAvailable(1) == 0 {
		return false
	}
	return true
}

func CreateAPI(c *API, db storage.Storage) (id uint64, err error) {
	id, err = db.Create(global.Config.ApiPrefix, c)
	return
}

func GetAPI(id uint64, db storage.Storage) (c *API, err error) {
	c = &API{}
	err = db.Get(global.Config.ApiPrefix, id, c)
	return
}

func WalkAPIs(start uint64, limit int64, walking func(e storage.Element) error, db storage.Storage) (nextID uint64, err error) {
	nextID, err = db.Walk(global.Config.ApiPrefix, start, limit, func() storage.Element {
		return &API{}
	}, walking)
	return
}

func UpdateAPI(c *API, db storage.Storage) (err error) {
	err = db.Update(global.Config.ApiPrefix, c)
	return
}

func DeleteAPI(id uint64, db storage.Storage) (err error) {
	err = db.Delete(global.Config.ApiPrefix, fmt.Sprintf("%d", id))
	return
}
