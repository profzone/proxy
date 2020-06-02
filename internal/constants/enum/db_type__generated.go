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
		"ETCD":    "etcd",
		"MONGODB": "mongodb",
		"MYSQL":   "mysql",
	})
}

func ParseDbTypeFromString(s string) (DbType, error) {
	switch s {
	case "":
		return DB_TYPE_UNKNOWN, nil
	case "ETCD":
		return DB_TYPE__ETCD, nil
	case "MONGODB":
		return DB_TYPE__MONGODB, nil
	case "MYSQL":
		return DB_TYPE__MYSQL, nil
	}
	return DB_TYPE_UNKNOWN, InvalidDbType
}

func ParseDbTypeFromLabelString(s string) (DbType, error) {
	switch s {
	case "":
		return DB_TYPE_UNKNOWN, nil
	case "etcd":
		return DB_TYPE__ETCD, nil
	case "mongodb":
		return DB_TYPE__MONGODB, nil
	case "mysql":
		return DB_TYPE__MYSQL, nil
	}
	return DB_TYPE_UNKNOWN, InvalidDbType
}

func (DbType) EnumType() string {
	return "DbType"
}

func (DbType) Enums() map[int][]string {
	return map[int][]string{
		int(DB_TYPE__ETCD):    {"ETCD", "etcd"},
		int(DB_TYPE__MONGODB): {"MONGODB", "mongodb"},
		int(DB_TYPE__MYSQL):   {"MYSQL", "mysql"},
	}
}

func (v DbType) String() string {
	switch v {
	case DB_TYPE_UNKNOWN:
		return ""
	case DB_TYPE__ETCD:
		return "ETCD"
	case DB_TYPE__MONGODB:
		return "MONGODB"
	case DB_TYPE__MYSQL:
		return "MYSQL"
	}
	return "UNKNOWN"
}

func (v DbType) Label() string {
	switch v {
	case DB_TYPE_UNKNOWN:
		return ""
	case DB_TYPE__ETCD:
		return "etcd"
	case DB_TYPE__MONGODB:
		return "mongodb"
	case DB_TYPE__MYSQL:
		return "mysql"
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
