package enum

//go:generate eden generate enum --type-name=DbType
// api:enum
type DbType uint8

// DB类型
const (
	DB_TYPE_UNKNOWN DbType = iota
	DB_TYPE__ETCD          // etcd
)
