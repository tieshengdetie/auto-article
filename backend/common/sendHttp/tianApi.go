package sendHttp

import (
	"AutoArticle/common/sendHttp/types"
	"AutoArticle/global"
	"AutoArticle/utils"
	"context"
	"fmt"
)

const TianApiType = "tianApi"

var TianApi = TianApiHttp{}

type TianApiHttp struct{}

var TianApiHttpUrlMap = urlMap{
	"nethot":      "/nethot/index",      // 百度热门搜索
	"weibohot":    "/weibohot/index",    // 微博热门搜索
	"douyinhot":   "/douyinhot/index",   // 抖音热门搜索
	"wxhottopic":  "/wxhottopic/index",  // 微信热门话题
	"toutiaohot":  "/toutiaohot/index",  // 今日头条热搜
	"allnews":     "/allnews/index",     // 指定分类新闻
	"generalnews": "/generalnews/index", // 综合新闻
}

func init() {
	registerUrl(TianApiType, TianApiHttpUrlMap)
}

// Config 配置请求参数
//
//	@Description: 配置请求头、完整URL
//	@receiver t
//	@param ctx
//	@param urlKey
//	@return headers
//	@return fullUrl
func (t *TianApiHttp) Config(ctx context.Context, urlKey string) (headers map[string]string, fullUrl string) {
	headers = make(map[string]string)
	config := global.ServerConfig.TianApiConfig
	fullUrl = config.BaseUrl + urlMapping[TianApiType][urlKey] + "?key=" + config.AppKey
	return
}

// Request 发送HTTP请求
//
//	@Description: 发送HTTP请求并解析响应
//	@receiver t
//	@param ctx
//	@param urlKey
//	@param receiver
//	@return err
func (t *TianApiHttp) Request(ctx context.Context, urlKey string, receiver interface{}) (err error) {
	headers, fullUrl := t.Config(ctx, urlKey)
	return utils.HTTPRequestByReceiver(utils.GET, fullUrl, nil, headers, nil, receiver)
}

// RequestWithParams 带参数的请求方法
//
//	@Description: 发送带查询参数的HTTP请求
//	@receiver t
//	@param ctx
//	@param urlKey
//	@param params
//	@param receiver
//	@return err
func (t *TianApiHttp) RequestWithParams(ctx context.Context, urlKey string, params map[string]interface{}, receiver interface{}) (err error) {
	headers, fullUrl := t.Config(ctx, urlKey)
	return utils.HTTPRequestByReceiver(utils.GET, fullUrl, params, headers, nil, receiver)
}

// GetNetHot 获取百度热门搜索
//
//	@Description: 获取百度热门搜索榜单
//	@receiver t
//	@param ctx
//	@return resp
//	@return err
func (t *TianApiHttp) GetNetHot(ctx context.Context) (resp *types.NetHotResp, err error) {
	resp = &types.NetHotResp{}
	err = t.Request(ctx, "nethot", resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("获取百度热搜失败: %s", resp.Msg)
	}
	return resp, nil
}

// GetWeiboHot 获取微博热门搜索
//
//	@Description: 获取微博热门搜索榜单
//	@receiver t
//	@param ctx
//	@return resp
//	@return err
func (t *TianApiHttp) GetWeiboHot(ctx context.Context) (resp *types.WeiboHotResp, err error) {
	resp = &types.WeiboHotResp{}
	err = t.Request(ctx, "weibohot", resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("获取微博热搜失败: %s", resp.Msg)
	}
	return resp, nil
}

// GetDouyinHot 获取抖音热门搜索
//
//	@Description: 获取抖音热门搜索榜单
//	@receiver t
//	@param ctx
//	@return resp
//	@return err
func (t *TianApiHttp) GetDouyinHot(ctx context.Context) (resp *types.DouyinHotResp, err error) {
	resp = &types.DouyinHotResp{}
	err = t.Request(ctx, "douyinhot", resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("获取抖音热搜失败: %s", resp.Msg)
	}
	return resp, nil
}

// GetWxHotTopic 获取微信热门话题
//
//	@Description: 获取微信热门话题榜单
//	@receiver t
//	@param ctx
//	@return resp
//	@return err
func (t *TianApiHttp) GetWxHotTopic(ctx context.Context) (resp *types.WxHotTopicResp, err error) {
	resp = &types.WxHotTopicResp{}
	err = t.Request(ctx, "wxhottopic", resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("获取微信热榜失败: %s", resp.Msg)
	}
	return resp, nil
}

// GetToutiaoHot 获取今日头条热搜
//
//	@Description: 获取今日头条热搜榜单
//	@receiver t
//	@param ctx
//	@return resp
//	@return err
func (t *TianApiHttp) GetToutiaoHot(ctx context.Context) (resp *types.ToutiaoHotResp, err error) {
	resp = &types.ToutiaoHotResp{}
	err = t.Request(ctx, "toutiaohot", resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("获取今日头条热搜失败: %s", resp.Msg)
	}
	return resp, nil
}

// GetAllNews 获取指定分类新闻
//
//	@Description: 获取指定分类的新闻列表，支持关键词搜索和翻页
//	@receiver t
//	@param ctx
//	@param col 频道ID（必填）
//	@param num 返回数量（1-50，默认10）
//	@param page 页码（默认1）
//	@param rand 是否随机获取（0不随机，1随机）
//	@param word 搜索关键词（可选）
//	@return resp
//	@return err
func (t *TianApiHttp) GetAllNews(ctx context.Context, col, num, page, rand int, word string) (resp *types.AllNewsResp, err error) {
	resp = &types.AllNewsResp{}
	params := map[string]interface{}{
		"col":  col,
		"num":  num,
		"page": page,
		"rand": rand,
		"form": 1,
	}
	if word != "" {
		params["word"] = word
	}
	err = t.RequestWithParams(ctx, "allnews", params, resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("获取新闻失败: %s", resp.Msg)
	}
	return resp, nil
}

// GetGeneralNews 获取综合新闻
//
//	@Description: 获取综合新闻列表，支持关键词搜索、翻页和指定来源
//	@receiver t
//	@param ctx
//	@param num 返回数量（1-50，默认10）
//	@param page 页码（默认0）
//	@param rand 是否随机获取（0不随机，1随机）
//	@param word 搜索关键词（可选）
//	@param source 指定来源（可选）
//	@return resp
//	@return err
func (t *TianApiHttp) GetGeneralNews(ctx context.Context, num, page, rand int, word, source string) (resp *types.GeneralNewsResp, err error) {
	resp = &types.GeneralNewsResp{}
	params := map[string]interface{}{
		"num":  num,
		"page": page,
		"rand": rand,
		"form": 1,
	}
	if word != "" {
		params["word"] = word
	}
	if source != "" {
		params["source"] = source
	}
	err = t.RequestWithParams(ctx, "generalnews", params, resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("获取综合新闻失败: %s", resp.Msg)
	}
	return resp, nil
}
