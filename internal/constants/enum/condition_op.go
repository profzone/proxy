package enum

//go:generate eden generate enum --type-name=ConditionOp
// api:enum
type ConditionOp uint8

// 条件操作
const (
	CONDITION_OP_UNKNOWN ConditionOp = iota
	CONDITION_OP__EQ                 // =
	CONDITION_OP__NOT_EQ             // !=
	CONDITION_OP__LT                 // <
	CONDITION_OP__LTE                // <=
	CONDITION_OP__GT                 // >
	CONDITION_OP__GTE                // >=
)
