package enum

import (
	"bytes"
	"encoding"
	"errors"

	github_com_profzone_eden_framework_pkg_enumeration "github.com/profzone/eden-framework/pkg/enumeration"
)

var InvalidLoadBalanceType = errors.New("invalid LoadBalanceType")

func init() {
	github_com_profzone_eden_framework_pkg_enumeration.RegisterEnums("LoadBalanceType", map[string]string{
		"ROUND_ROBIN": "RoundRobin 算法",
	})
}

func ParseLoadBalanceTypeFromString(s string) (LoadBalanceType, error) {
	switch s {
	case "":
		return LOAD_BALANCE_TYPE_UNKNOWN, nil
	case "ROUND_ROBIN":
		return LOAD_BALANCE_TYPE__ROUND_ROBIN, nil
	}
	return LOAD_BALANCE_TYPE_UNKNOWN, InvalidLoadBalanceType
}

func ParseLoadBalanceTypeFromLabelString(s string) (LoadBalanceType, error) {
	switch s {
	case "":
		return LOAD_BALANCE_TYPE_UNKNOWN, nil
	case "RoundRobin 算法":
		return LOAD_BALANCE_TYPE__ROUND_ROBIN, nil
	}
	return LOAD_BALANCE_TYPE_UNKNOWN, InvalidLoadBalanceType
}

func (LoadBalanceType) EnumType() string {
	return "LoadBalanceType"
}

func (LoadBalanceType) Enums() map[int][]string {
	return map[int][]string{
		int(LOAD_BALANCE_TYPE__ROUND_ROBIN): {"ROUND_ROBIN", "RoundRobin 算法"},
	}
}

func (v LoadBalanceType) String() string {
	switch v {
	case LOAD_BALANCE_TYPE_UNKNOWN:
		return ""
	case LOAD_BALANCE_TYPE__ROUND_ROBIN:
		return "ROUND_ROBIN"
	}
	return "UNKNOWN"
}

func (v LoadBalanceType) Label() string {
	switch v {
	case LOAD_BALANCE_TYPE_UNKNOWN:
		return ""
	case LOAD_BALANCE_TYPE__ROUND_ROBIN:
		return "RoundRobin 算法"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*LoadBalanceType)(nil)

func (v LoadBalanceType) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidLoadBalanceType
	}
	return []byte(str), nil
}

func (v *LoadBalanceType) UnmarshalText(data []byte) (err error) {
	*v, err = ParseLoadBalanceTypeFromString(string(bytes.ToUpper(data)))
	return
}
