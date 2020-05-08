package enum

import (
	"bytes"
	"encoding"
	"errors"

	github_com_profzone_eden_framework_pkg_enumeration "github.com/profzone/eden-framework/pkg/enumeration"
)

var InvalidApiStatus = errors.New("invalid ApiStatus")

func init() {
	github_com_profzone_eden_framework_pkg_enumeration.RegisterEnums("ApiStatus", map[string]string{
		"DOWN": "下线",
		"UP":   "启用",
	})
}

func ParseApiStatusFromString(s string) (ApiStatus, error) {
	switch s {
	case "":
		return API_STATUS_UNKNOWN, nil
	case "DOWN":
		return API_STATUS__DOWN, nil
	case "UP":
		return API_STATUS__UP, nil
	}
	return API_STATUS_UNKNOWN, InvalidApiStatus
}

func ParseApiStatusFromLabelString(s string) (ApiStatus, error) {
	switch s {
	case "":
		return API_STATUS_UNKNOWN, nil
	case "下线":
		return API_STATUS__DOWN, nil
	case "启用":
		return API_STATUS__UP, nil
	}
	return API_STATUS_UNKNOWN, InvalidApiStatus
}

func (ApiStatus) EnumType() string {
	return "ApiStatus"
}

func (ApiStatus) Enums() map[int][]string {
	return map[int][]string{
		int(API_STATUS__DOWN): {"DOWN", "下线"},
		int(API_STATUS__UP):   {"UP", "启用"},
	}
}

func (v ApiStatus) String() string {
	switch v {
	case API_STATUS_UNKNOWN:
		return ""
	case API_STATUS__DOWN:
		return "DOWN"
	case API_STATUS__UP:
		return "UP"
	}
	return "UNKNOWN"
}

func (v ApiStatus) Label() string {
	switch v {
	case API_STATUS_UNKNOWN:
		return ""
	case API_STATUS__DOWN:
		return "下线"
	case API_STATUS__UP:
		return "启用"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*ApiStatus)(nil)

func (v ApiStatus) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidApiStatus
	}
	return []byte(str), nil
}

func (v *ApiStatus) UnmarshalText(data []byte) (err error) {
	*v, err = ParseApiStatusFromString(string(bytes.ToUpper(data)))
	return
}
