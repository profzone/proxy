package enum

//go:generate eden generate enum --type-name=ServerType
// api:enum
type ServerType uint8

// 服务器类型
const (
	SERVER_TYPE_UNKNOWN      ServerType = iota
	SERVER_TYPE__WEB_SERVICE            // Web服务
	SERVER_TYPE__DATABASE               // 数据库服务
)
