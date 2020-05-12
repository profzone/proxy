package modules

import (
	"fmt"
	str "github.com/profzone/eden-framework/pkg/strings"
	"github.com/valyala/fasthttp"
	"longhorn/proxy/pkg/route"
	"strings"
)

const (
	conditionOpEqual    = "="
	conditionOpNotEqual = "!="
	conditionOpGte      = ">="
	conditionOpGt       = ">"
	conditionOpLte      = "<="
	conditionOpLt       = "<"
	conditionOpContain  = "~"
)

var opArray = []string{
	conditionOpEqual,
	conditionOpNotEqual,
	conditionOpGte,
	conditionOpGt,
	conditionOpLte,
	conditionOpLt,
	conditionOpContain,
}

func getOpFuncByStr(opStr string) opFunc {
	switch opStr {
	case conditionOpEqual:
		return opEqual
	case conditionOpNotEqual:
		return opNotEqual
	case conditionOpGte:
		return opGte
	case conditionOpGt:
		return opGt
	case conditionOpLte:
		return opLte
	case conditionOpLt:
		return opLt
	case conditionOpContain:
		return opContain
	default:
		return nil
	}
}

type routerCondition struct {
	origin string
	rules  []rule
	params map[string]string
}

func newRouterCondition(condition string, params route.Params) *routerCondition {
	if condition == "" {
		return nil
	}
	c := &routerCondition{
		origin: condition,
		rules:  make([]rule, 0),
		params: make(map[string]string),
	}

	for _, p := range params {
		c.params[p.Key] = p.Value
	}

	if err := c.scan(); err != nil {
		return nil
	}

	return c
}

func (c *routerCondition) scan() error {
	rulesStr := strings.Split(c.origin, ",")
	for _, ruleStr := range rulesStr {
		rule, err := c.scanRule(ruleStr)
		if err != nil {
			return err
		}
		c.rules = append(c.rules, rule)
	}
	return nil
}

func (c *routerCondition) scanRule(rule string) (rule, error) {
	index, indexInSlice := str.StringIndexInSlice(rule, opArray)
	if index < 0 {
		return nil, fmt.Errorf("syntax error: no valid op command found")
	}
	key := rule[:index]
	opStr := opArray[indexInSlice]
	value := rule[index+len(opStr):]
	opFunc := getOpFuncByStr(opStr)
	if opFunc == nil {
		return nil, fmt.Errorf("syntax error: no valid op command found")
	}
	return newRule(key, value, opFunc)
}

func (c *routerCondition) Match(req *fasthttp.Request) bool {
	for _, r := range c.rules {
		if !r.match(req, c.params) {
			return false
		}
	}
	return true
}

func (c *routerCondition) UnmarshalJSON([]byte) error {
	panic("implement me")
}

func (c routerCondition) MarshalJSON() ([]byte, error) {
	panic("implement me")
}

func (c *routerCondition) GobDecode([]byte) error {
	panic("implement me")
}

func (c routerCondition) GobEncode() ([]byte, error) {
	panic("implement me")
}
