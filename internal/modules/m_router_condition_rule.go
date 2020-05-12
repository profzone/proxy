package modules

import (
	"github.com/buger/jsonparser"
	"github.com/valyala/fasthttp"
	"strings"
)

type opFunc func(origin string, value string) bool

func opEqual(origin string, value string) bool {
	if strings.Compare(origin, value) == 0 {
		return true
	}
	return false
}

func opNotEqual(origin string, value string) bool {
	if strings.Compare(origin, value) != 0 {
		return true
	}
	return false
}

func opContain(origin string, value string) bool {
	if strings.Contains(origin, value) {
		return true
	}
	return false
}

func opGte(origin string, value string) bool {
	flag := strings.Compare(origin, value)
	if flag >= 0 {
		return true
	}
	return false
}

func opGt(origin string, value string) bool {
	flag := strings.Compare(origin, value)
	if flag > 0 {
		return true
	}
	return false
}

func opLte(origin string, value string) bool {
	flag := strings.Compare(origin, value)
	if flag <= 0 {
		return true
	}
	return false
}

func opLt(origin string, value string) bool {
	flag := strings.Compare(origin, value)
	if flag < 0 {
		return true
	}
	return false
}

type rule interface {
	match(req *fasthttp.Request) bool
}

type bodyRule struct {
	key   string
	op    opFunc
	value string
}

func (b bodyRule) match(req *fasthttp.Request) bool {
	path := strings.Split(b.key, ".")
	origin, err := jsonparser.GetString(req.Body(), path...)
	if err != nil {
		return false
	}
	return b.op(origin, b.value)
}

type headerRule struct {
	key   string
	op    opFunc
	value string
}

func (b headerRule) match(req *fasthttp.Request) bool {
	origin := string(req.Header.Peek(b.key))
	return b.op(origin, b.value)
}

type queryRule struct {
	key   string
	op    opFunc
	value string
}

func (b queryRule) match(req *fasthttp.Request) bool {
	origin := string(req.URI().QueryArgs().Peek(b.key))
	return b.op(origin, b.value)
}

type cookieRule struct {
	key   string
	op    opFunc
	value string
}

func (b cookieRule) match(req *fasthttp.Request) bool {
	origin := string(req.Header.Cookie(b.key))
	return b.op(origin, b.value)
}
