package http

type RequestPagination struct {
	// 分页偏移量
	Offset int32 `json:"offset" name:"offset" in:"query" default:"10" validate:"@int32[0,100]"`
	// 每页数据量
	Size int32 `json:"size" name:"size" in:"query" default:"10" validate:"@int32[10,100]"`
}

type IDResponse struct {
	// ID
	ID uint64 `json:"id"`
}
