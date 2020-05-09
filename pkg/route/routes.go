package route

import "sync"

func countParams(path string) uint16 {
	var n uint
	for i := range []byte(path) {
		switch path[i] {
		case ':', '*':
			n++
		}
	}
	return uint16(n)
}

type Routes struct {
	trees      map[string]*node
	paramsPool sync.Pool
	maxParams  uint16
}

func NewRoutes() *Routes {
	return &Routes{
		trees: make(map[string]*node),
	}
}

func (r *Routes) getParams() *Params {
	ps := r.paramsPool.Get().(*Params)
	*ps = (*ps)[0:0] // reset slice
	return ps
}

func (r *Routes) putParams(ps *Params) {
	if ps != nil {
		r.paramsPool.Put(ps)
	}
}

func (r *Routes) Handle(method, path string, apiID uint64) error {
	varsCount := uint16(0)

	if method == "" {
		panic("method must not be empty")
	}
	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}
	if apiID == 0 {
		panic("API must be set")
	}

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root
	}

	err := root.addRoute(path, apiID)
	if err != nil {
		return err
	}

	// Update maxParams
	if paramsCount := countParams(path); paramsCount+varsCount > r.maxParams {
		r.maxParams = paramsCount + varsCount
	}

	// Lazy-init paramsPool alloc func
	if r.paramsPool.New == nil && r.maxParams > 0 {
		r.paramsPool.New = func() interface{} {
			ps := make(Params, 0, r.maxParams)
			return &ps
		}
	}

	return nil
}

func (r *Routes) Lookup(method, path string) (uint64, Params, bool) {
	if root := r.trees[method]; root != nil {
		apiID, ps, tsr := root.getValue(path, r.getParams)
		if apiID == 0 {
			r.putParams(ps)
			return 0, nil, tsr
		}
		if ps == nil {
			return apiID, nil, tsr
		}
		return apiID, *ps, tsr
	}
	return 0, nil, false
}
