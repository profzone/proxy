package modules

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

const (
	conditionTokenDot    = "."
	conditionExprContext = "context"
	conditionExprPath    = "path"
	conditionExprQuery   = "query"
	conditionExprCookie  = "cookie"
	conditionExprHeader  = "header"
	conditionExprBody    = "body"
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

func opStrGte(origin string, value string) bool {
	flag := strings.Compare(origin, value)
	if flag >= 0 {
		return true
	}
	return false
}

func opStrGt(origin string, value string) bool {
	flag := strings.Compare(origin, value)
	if flag > 0 {
		return true
	}
	return false
}

func opStrLte(origin string, value string) bool {
	flag := strings.Compare(origin, value)
	if flag <= 0 {
		return true
	}
	return false
}

func opStrLt(origin string, value string) bool {
	flag := strings.Compare(origin, value)
	if flag < 0 {
		return true
	}
	return false
}

func opGte(origin string, value string) bool {
	originDigit, _ := strconv.Atoi(origin)
	valueDigit, _ := strconv.Atoi(value)
	if originDigit >= valueDigit {
		return true
	}
	return false
}

func opGt(origin string, value string) bool {
	originDigit, _ := strconv.Atoi(origin)
	valueDigit, _ := strconv.Atoi(value)
	if originDigit > valueDigit {
		return true
	}
	return false
}

func opLte(origin string, value string) bool {
	originDigit, _ := strconv.Atoi(origin)
	valueDigit, _ := strconv.Atoi(value)
	if originDigit <= valueDigit {
		return true
	}
	return false
}

func opLt(origin string, value string) bool {
	originDigit, _ := strconv.Atoi(origin)
	valueDigit, _ := strconv.Atoi(value)
	if originDigit < valueDigit {
		return true
	}
	return false
}

type rule interface {
	match(req *fasthttp.Request, params map[string]string) bool
}

func newRule(key, value string, op opFunc) (rule, error) {
	keyArray := strings.Split(key, conditionTokenDot)
	if len(keyArray) < 2 {
		return nil, fmt.Errorf("syntax error: the key must have at least 2 segments")
	}
	switch keyArray[0] {
	case conditionExprBody:
		return bodyRule{
			key:   keyArray[1:],
			op:    op,
			value: value,
		}, nil
	case conditionExprHeader:
		return headerRule{
			key:   keyArray[1],
			op:    op,
			value: value,
		}, nil
	case conditionExprQuery:
		return queryRule{
			key:   keyArray[1],
			op:    op,
			value: value,
		}, nil
	case conditionExprCookie:
		return cookieRule{
			key:   keyArray[1],
			op:    op,
			value: value,
		}, nil
	case conditionExprPath:
		return pathRule{
			key:   keyArray[1],
			op:    op,
			value: value,
		}, nil
	default:
		return nil, fmt.Errorf("syntax error: invalid expression %s", keyArray[0])
	}
}

type bodyRule struct {
	key   []string
	op    opFunc
	value string
}

func (b bodyRule) match(req *fasthttp.Request, params map[string]string) bool {
	origin, err := jsonparser.GetString(req.Body(), b.key...)
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

func (b headerRule) match(req *fasthttp.Request, params map[string]string) bool {
	origin := string(req.Header.Peek(b.key))
	return b.op(origin, b.value)
}

type queryRule struct {
	key   string
	op    opFunc
	value string
}

func (b queryRule) match(req *fasthttp.Request, params map[string]string) bool {
	origin := string(req.URI().QueryArgs().Peek(b.key))
	return b.op(origin, b.value)
}

type cookieRule struct {
	key   string
	op    opFunc
	value string
}

func (b cookieRule) match(req *fasthttp.Request, params map[string]string) bool {
	origin := string(req.Header.Cookie(b.key))
	return b.op(origin, b.value)
}

type pathRule struct {
	key   string
	op    opFunc
	value string
}

func (b pathRule) match(req *fasthttp.Request, params map[string]string) bool {
	origin := params[b.key]
	return b.op(origin, b.value)
}
