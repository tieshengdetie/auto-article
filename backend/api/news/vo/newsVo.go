package vo

// NewsItem 新闻条目
type NewsItem struct {
	ID          string `json:"id"`
	Url         string `json:"url"`
	Ctime       string `json:"ctime"`
	Title       string `json:"title"`
	PicUrl      string `json:"picUrl"`
	Source      string `json:"source"`
	Description string `json:"description"`
}

// NewsResp 新闻响应
type NewsResp struct {
	List    []NewsItem `json:"list"`
	Allnum  int        `json:"allnum"`
	Curpage int        `json:"curpage"`
}

// NewsChannel 新闻频道
type NewsChannel struct {
	ChannelID   int    `json:"channelId"`
	ChannelName string `json:"channelName"`
	Description string `json:"description"`
}
