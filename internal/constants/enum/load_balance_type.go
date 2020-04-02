package enum

//go:generate eden generate enum --type-name=LoadBalanceType
// api:enum
type LoadBalanceType uint8

// 负载均衡类型
const (
	LOAD_BALANCE_TYPE_UNKNOWN      LoadBalanceType = iota
	LOAD_BALANCE_TYPE__ROUND_ROBIN                 // RoundRobin 算法
)
