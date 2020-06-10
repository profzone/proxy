package storage

import "longhorn/proxy/internal/constants/enum"

type Condition struct {
	Key string
	Op  enum.ConditionOp
	Val interface{}
}

func WithConditionKey(key string) *Condition {
	return &Condition{
		Key: key,
	}
}

func (c *Condition) Eq(val interface{}) *Condition {
	c.Op = enum.CONDITION_OP__EQ
	c.Val = val
	return c
}
