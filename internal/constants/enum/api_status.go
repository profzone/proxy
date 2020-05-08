package enum

//go:generate eden generate enum --type-name=ApiStatus
// api:enum
type ApiStatus uint8

// API状态
const (
	API_STATUS_UNKNOWN ApiStatus = iota
	API_STATUS__UP               // 启用
	API_STATUS__DOWN             // 下线
)
