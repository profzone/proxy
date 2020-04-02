package enum

import (
	"bytes"
	"encoding"
	"errors"

	github_com_profzone_eden_framework_pkg_enumeration "github.com/profzone/eden-framework/pkg/enumeration"
)

var InvalidDbType = errors.New("invalid DbType")

func init() {
	github_com_profzone_eden_framework_pkg_enumeration.RegisterEnums("DbType", map[string]string{
		"ETCD": "etcd",
	})
}

func ParseDbTypeFromString(s string) (DbType, error) {
	switch s {
	case "":
		return DB_TYPE_UNKNOWN, nil
	case "ETCD":
		return DB_TYPE__ETCD, nil
	}
	return DB_TYPE_UNKNOWN, InvalidDbType
}

func ParseDbTypeFromLabelString(s string) (DbType, error) {
	switch s {
	case "":
		return DB_TYPE_UNKNOWN, nil
	case "etcd":
		return DB_TYPE__ETCD, nil
	}
	return DB_TYPE_UNKNOWN, InvalidDbType
}

func (DbType) EnumType() string {
	return "DbType"
}

func (DbType) Enums() map[int][]string {
	return map[int][]string{
		int(DB_TYPE__ETCD): {"ETCD", "etcd"},
	}
}

func (v DbType) String() string {
	switch v {
	case DB_TYPE_UNKNOWN:
		return ""
	case DB_TYPE__ETCD:
		return "ETCD"
	}
	return "UNKNOWN"
}

func (v DbType) Label() string {
	switch v {
	case DB_TYPE_UNKNOWN:
		return ""
	case DB_TYPE__ETCD:
		return "etcd"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*DbType)(nil)

func (v DbType) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidDbType
	}
	return []byte(str), nil
}

func (v *DbType) UnmarshalText(data []byte) (err error) {
	*v, err = ParseDbTypeFromString(string(bytes.ToUpper(data)))
	return
}
