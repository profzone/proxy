package enum

import (
	"bytes"
	"encoding"
	"errors"

	github_com_profzone_eden_framework_pkg_enumeration "github.com/profzone/eden-framework/pkg/enumeration"
)

var InvalidServerType = errors.New("invalid ServerType")

func init() {
	github_com_profzone_eden_framework_pkg_enumeration.RegisterEnums("ServerType", map[string]string{
		"DATABASE":    "数据库服务",
		"WEB_SERVICE": "Web服务",
	})
}

func ParseServerTypeFromString(s string) (ServerType, error) {
	switch s {
	case "":
		return SERVER_TYPE_UNKNOWN, nil
	case "DATABASE":
		return SERVER_TYPE__DATABASE, nil
	case "WEB_SERVICE":
		return SERVER_TYPE__WEB_SERVICE, nil
	}
	return SERVER_TYPE_UNKNOWN, InvalidServerType
}

func ParseServerTypeFromLabelString(s string) (ServerType, error) {
	switch s {
	case "":
		return SERVER_TYPE_UNKNOWN, nil
	case "数据库服务":
		return SERVER_TYPE__DATABASE, nil
	case "Web服务":
		return SERVER_TYPE__WEB_SERVICE, nil
	}
	return SERVER_TYPE_UNKNOWN, InvalidServerType
}

func (ServerType) EnumType() string {
	return "ServerType"
}

func (ServerType) Enums() map[int][]string {
	return map[int][]string{
		int(SERVER_TYPE__DATABASE):    {"DATABASE", "数据库服务"},
		int(SERVER_TYPE__WEB_SERVICE): {"WEB_SERVICE", "Web服务"},
	}
}

func (v ServerType) String() string {
	switch v {
	case SERVER_TYPE_UNKNOWN:
		return ""
	case SERVER_TYPE__DATABASE:
		return "DATABASE"
	case SERVER_TYPE__WEB_SERVICE:
		return "WEB_SERVICE"
	}
	return "UNKNOWN"
}

func (v ServerType) Label() string {
	switch v {
	case SERVER_TYPE_UNKNOWN:
		return ""
	case SERVER_TYPE__DATABASE:
		return "数据库服务"
	case SERVER_TYPE__WEB_SERVICE:
		return "Web服务"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*ServerType)(nil)

func (v ServerType) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidServerType
	}
	return []byte(str), nil
}

func (v *ServerType) UnmarshalText(data []byte) (err error) {
	*v, err = ParseServerTypeFromString(string(bytes.ToUpper(data)))
	return
}
