package enum

import (
	"bytes"
	"encoding"
	"errors"

	github_com_profzone_eden_framework_pkg_enumeration "github.com/profzone/eden-framework/pkg/enumeration"
)

var InvalidConditionOp = errors.New("invalid ConditionOp")

func init() {
	github_com_profzone_eden_framework_pkg_enumeration.RegisterEnums("ConditionOp", map[string]string{
		"EQ":     "=",
		"GT":     ">",
		"GTE":    ">=",
		"LT":     "<",
		"LTE":    "<=",
		"NOT_EQ": "!=",
	})
}

func ParseConditionOpFromString(s string) (ConditionOp, error) {
	switch s {
	case "":
		return CONDITION_OP_UNKNOWN, nil
	case "EQ":
		return CONDITION_OP__EQ, nil
	case "GT":
		return CONDITION_OP__GT, nil
	case "GTE":
		return CONDITION_OP__GTE, nil
	case "LT":
		return CONDITION_OP__LT, nil
	case "LTE":
		return CONDITION_OP__LTE, nil
	case "NOT_EQ":
		return CONDITION_OP__NOT_EQ, nil
	}
	return CONDITION_OP_UNKNOWN, InvalidConditionOp
}

func ParseConditionOpFromLabelString(s string) (ConditionOp, error) {
	switch s {
	case "":
		return CONDITION_OP_UNKNOWN, nil
	case "=":
		return CONDITION_OP__EQ, nil
	case ">":
		return CONDITION_OP__GT, nil
	case ">=":
		return CONDITION_OP__GTE, nil
	case "<":
		return CONDITION_OP__LT, nil
	case "<=":
		return CONDITION_OP__LTE, nil
	case "!=":
		return CONDITION_OP__NOT_EQ, nil
	}
	return CONDITION_OP_UNKNOWN, InvalidConditionOp
}

func (ConditionOp) EnumType() string {
	return "ConditionOp"
}

func (ConditionOp) Enums() map[int][]string {
	return map[int][]string{
		int(CONDITION_OP__EQ):     {"EQ", "="},
		int(CONDITION_OP__GT):     {"GT", ">"},
		int(CONDITION_OP__GTE):    {"GTE", ">="},
		int(CONDITION_OP__LT):     {"LT", "<"},
		int(CONDITION_OP__LTE):    {"LTE", "<="},
		int(CONDITION_OP__NOT_EQ): {"NOT_EQ", "!="},
	}
}

func (v ConditionOp) String() string {
	switch v {
	case CONDITION_OP_UNKNOWN:
		return ""
	case CONDITION_OP__EQ:
		return "EQ"
	case CONDITION_OP__GT:
		return "GT"
	case CONDITION_OP__GTE:
		return "GTE"
	case CONDITION_OP__LT:
		return "LT"
	case CONDITION_OP__LTE:
		return "LTE"
	case CONDITION_OP__NOT_EQ:
		return "NOT_EQ"
	}
	return "UNKNOWN"
}

func (v ConditionOp) Label() string {
	switch v {
	case CONDITION_OP_UNKNOWN:
		return ""
	case CONDITION_OP__EQ:
		return "="
	case CONDITION_OP__GT:
		return ">"
	case CONDITION_OP__GTE:
		return ">="
	case CONDITION_OP__LT:
		return "<"
	case CONDITION_OP__LTE:
		return "<="
	case CONDITION_OP__NOT_EQ:
		return "!="
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*ConditionOp)(nil)

func (v ConditionOp) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidConditionOp
	}
	return []byte(str), nil
}

func (v *ConditionOp) UnmarshalText(data []byte) (err error) {
	*v, err = ParseConditionOpFromString(string(bytes.ToUpper(data)))
	return
}
