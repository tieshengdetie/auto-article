package types

// TianApiBaseResp 天聚数行基础响应结构
type TianApiBaseResp struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

// ==================== 百度热搜 ====================

// NetHotResp 百度热门搜索响应
type NetHotResp struct {
	Code   int          `json:"code"`
	Msg    string       `json:"msg"`
	Result NetHotResult `json:"result"`
}

// NetHotResult 百度热搜结果
type NetHotResult struct {
	List []NetHotItem `json:"list"`
}

// NetHotItem 百度热搜条目
type NetHotItem struct {
	Brief   string `json:"brief"`
	Index   string `json:"index"`
	Trend   string `json:"trend"`
	Keyword string `json:"keyword"`
}

// ==================== 微博热搜 ====================

// WeiboHotResp 微博热门搜索响应
type WeiboHotResp struct {
	Code   int            `json:"code"`
	Msg    string         `json:"msg"`
	Result WeiboHotResult `json:"result"`
}

// WeiboHotResult 微博热搜结果
type WeiboHotResult struct {
	List []WeiboHotItem `json:"list"`
}

// WeiboHotItem 微博热搜条目
type WeiboHotItem struct {
	Hottag     string `json:"hottag"`
	Hotword    string `json:"hotword"`
	Hotwordnum string `json:"hotwordnum"`
}

// ==================== 抖音热搜 ====================

// DouyinHotResp 抖音热门搜索响应
type DouyinHotResp struct {
	Code   int             `json:"code"`
	Msg    string          `json:"msg"`
	Result DouyinHotResult `json:"result"`
}

// DouyinHotResult 抖音热搜结果
type DouyinHotResult struct {
	List []DouyinHotItem `json:"list"`
}

// DouyinHotItem 抖音热搜条目
type DouyinHotItem struct {
	Word     string `json:"word"`
	Label    int    `json:"label"`
	Hotindex int    `json:"hotindex"`
}

// ==================== 微信热榜 ====================

// WxHotTopicResp 微信热门话题响应
type WxHotTopicResp struct {
	Code   int              `json:"code"`
	Msg    string           `json:"msg"`
	Result WxHotTopicResult `json:"result"`
}

// WxHotTopicResult 微信热榜结果
type WxHotTopicResult struct {
	List []WxHotTopicItem `json:"list"`
}

// WxHotTopicItem 微信热榜条目
type WxHotTopicItem struct {
	Word  string `json:"word"`
	Index int    `json:"index"`
}

// ==================== 今日头条热搜 ====================

// ToutiaoHotResp 今日头条热搜响应
type ToutiaoHotResp struct {
	Code   int              `json:"code"`
	Msg    string           `json:"msg"`
	Result ToutiaoHotResult `json:"result"`
}

// ToutiaoHotResult 今日头条热搜结果
type ToutiaoHotResult struct {
	List []ToutiaoHotItem `json:"list"`
}

// ToutiaoHotItem 今日头条热搜条目
type ToutiaoHotItem struct {
	Word     string `json:"word"`
	Hotindex int    `json:"hotindex"`
}

// ==================== 新闻相关 ====================

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

// NewsResult 新闻结果
type NewsResult struct {
	List    []NewsItem `json:"list"`
	Allnum  int        `json:"allnum"`
	Curpage int        `json:"curpage"`
}

// AllNewsResp 指定新闻分类响应
type AllNewsResp struct {
	Code   int        `json:"code"`
	Msg    string     `json:"msg"`
	Result NewsResult `json:"result"`
}

// GeneralNewsResp 综合新闻响应
type GeneralNewsResp struct {
	Code   int        `json:"code"`
	Msg    string     `json:"msg"`
	Result NewsResult `json:"result"`
}
