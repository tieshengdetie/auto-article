package dto

// HotSearchReq 热搜请求参数
type HotSearchReq struct {
	Type string `json:"type" binding:"required"`
}
