package modules

import (
	"github.com/juju/ratelimit"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"longhorn/proxy/internal/constants/enum"
	"longhorn/proxy/internal/models"
	"time"
)

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
	ipController *IPController
	// 最大QPS
	limiter *ratelimit.Bucket
	// TODO Validations
	// 反向代理调度
	Dispatchers []*Dispatcher `json:"dispatcher"`
	// TODO Fusion
}

func NewAPI(model *models.API) (*API, error) {
	var dispatchers = make([]*Dispatcher, 0)
	var ipController *IPController
	var limiter *ratelimit.Bucket
	var err error
	for _, d := range model.Dispatchers {
		dispatcher := NewDispatcher(&d)
		dispatchers = append(dispatchers, dispatcher)
	}

	if model.IPControl != "" {
		ipController, err = newIPController(model.IPControl)
		if err != nil {
			logrus.Errorf("[NewAPI] newIPController err: %v", err)
			return nil, err
		}
	}
	if model.MaxQPS > 0 {
		limiter = ratelimit.NewBucket(time.Second/time.Duration(model.MaxQPS), model.MaxQPS)
	}
	return &API{
		ID:           model.ID,
		Name:         model.Name,
		URLPattern:   model.URLPattern,
		Method:       model.Method,
		Status:       model.Status,
		ipController: ipController,
		limiter:      limiter,
		Dispatchers:  dispatchers,
	}, nil
}

func (v *API) SetIdentity(id uint64) {
	v.ID = id
}

func (v API) GetIdentity() uint64 {
	return v.ID
}

func (v *API) WalkDispatcher(walking func(dispatcher *Dispatcher) error) {
	for _, d := range v.Dispatchers {
		err := walking(d)
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
