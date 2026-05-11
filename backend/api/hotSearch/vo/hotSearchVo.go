package vo

// HotSearchItem 热搜条目
type HotSearchItem struct {
	// 通用字段
	Word     string `json:"word"`
	HotIndex int    `json:"hotIndex"`

	// 百度热搜字段
	Brief   string `json:"brief,omitempty"`
	Index   string `json:"index,omitempty"`
	Trend   string `json:"trend,omitempty"`
	Keyword string `json:"keyword,omitempty"`

	// 微博热搜字段
	HotTag     string `json:"hotTag,omitempty"`
	HotWord    string `json:"hotWord,omitempty"`
	HotWordNum string `json:"hotWordNum,omitempty"`

	// 抖音热搜字段
	Label int `json:"label,omitempty"`
}

// HotSearchResp 热搜响应
type HotSearchResp struct {
	Type string          `json:"type"`
	Name string          `json:"name"`
	List []HotSearchItem `json:"list"`
}
