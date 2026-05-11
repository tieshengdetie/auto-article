package service

import (
	"AutoArticle/api/hotSearch/dto"
	"AutoArticle/api/hotSearch/vo"
	"AutoArticle/common/sendHttp"
	"context"
	"fmt"
)

type HotSearchServer struct {
}

// 热搜类型映射
var hotSearchTypeMap = map[string]string{
	"nethot":     "百度热搜",
	"weibohot":   "微博热搜",
	"douyinhot":  "抖音热搜",
	"wxhottopic": "微信热榜",
	"toutiaohot": "今日头条",
}

// GetHotSearch 根据类型获取热搜数据
//
//	@Description: 根据传入的热搜类型，调用对应的接口获取数据
//	@receiver h
//	@param ctx
//	@param req
//	@return resp
//	@return err
func (h *HotSearchServer) GetHotSearch(ctx context.Context, req dto.HotSearchReq) (resp *vo.HotSearchResp, err error) {
	resp = &vo.HotSearchResp{
		Type: req.Type,
		Name: hotSearchTypeMap[req.Type],
	}

	switch req.Type {
	case "nethot":
		return h.getNetHot(ctx, resp)
	case "weibohot":
		return h.getWeiboHot(ctx, resp)
	case "douyinhot":
		return h.getDouyinHot(ctx, resp)
	case "wxhottopic":
		return h.getWxHotTopic(ctx, resp)
	case "toutiaohot":
		return h.getToutiaoHot(ctx, resp)
	default:
		return nil, fmt.Errorf("未知的热搜类型: %s", req.Type)
	}
}

// getNetHot 获取百度热搜
func (h *HotSearchServer) getNetHot(ctx context.Context, resp *vo.HotSearchResp) (*vo.HotSearchResp, error) {
	data, err := sendHttp.TianApi.GetNetHot(ctx)
	if err != nil {
		return nil, err
	}

	resp.List = make([]vo.HotSearchItem, 0, len(data.Result.List))
	for _, item := range data.Result.List {
		resp.List = append(resp.List, vo.HotSearchItem{
			Word:     item.Keyword,
			HotIndex: 0,
			Brief:    item.Brief,
			Index:    item.Index,
			Trend:    item.Trend,
			Keyword:  item.Keyword,
		})
	}

	return resp, nil
}

// getWeiboHot 获取微博热搜
func (h *HotSearchServer) getWeiboHot(ctx context.Context, resp *vo.HotSearchResp) (*vo.HotSearchResp, error) {
	data, err := sendHttp.TianApi.GetWeiboHot(ctx)
	if err != nil {
		return nil, err
	}

	resp.List = make([]vo.HotSearchItem, 0, len(data.Result.List))
	for _, item := range data.Result.List {
		resp.List = append(resp.List, vo.HotSearchItem{
			Word:       item.Hotword,
			HotIndex:   0,
			HotTag:     item.Hottag,
			HotWord:    item.Hotword,
			HotWordNum: item.Hotwordnum,
		})
	}

	return resp, nil
}

// getDouyinHot 获取抖音热搜
func (h *HotSearchServer) getDouyinHot(ctx context.Context, resp *vo.HotSearchResp) (*vo.HotSearchResp, error) {
	data, err := sendHttp.TianApi.GetDouyinHot(ctx)
	if err != nil {
		return nil, err
	}

	resp.List = make([]vo.HotSearchItem, 0, len(data.Result.List))
	for _, item := range data.Result.List {
		resp.List = append(resp.List, vo.HotSearchItem{
			Word:     item.Word,
			HotIndex: item.Hotindex,
			Label:    item.Label,
		})
	}

	return resp, nil
}

// getWxHotTopic 获取微信热榜
func (h *HotSearchServer) getWxHotTopic(ctx context.Context, resp *vo.HotSearchResp) (*vo.HotSearchResp, error) {
	data, err := sendHttp.TianApi.GetWxHotTopic(ctx)
	if err != nil {
		return nil, err
	}

	resp.List = make([]vo.HotSearchItem, 0, len(data.Result.List))
	for _, item := range data.Result.List {
		resp.List = append(resp.List, vo.HotSearchItem{
			Word:     item.Word,
			HotIndex: item.Index,
		})
	}

	return resp, nil
}

// getToutiaoHot 获取今日头条热搜
func (h *HotSearchServer) getToutiaoHot(ctx context.Context, resp *vo.HotSearchResp) (*vo.HotSearchResp, error) {
	data, err := sendHttp.TianApi.GetToutiaoHot(ctx)
	if err != nil {
		return nil, err
	}

	resp.List = make([]vo.HotSearchItem, 0, len(data.Result.List))
	for _, item := range data.Result.List {
		resp.List = append(resp.List, vo.HotSearchItem{
			Word:     item.Word,
			HotIndex: item.Hotindex,
		})
	}

	return resp, nil
}

// GetHotSearchTypes 获取所有支持的热搜类型
//
//	@Description: 获取系统支持的所有热搜类型列表
//	@receiver h
//	@return types
func (h *HotSearchServer) GetHotSearchTypes() []map[string]string {
	result := make([]map[string]string, 0, len(hotSearchTypeMap))
	for k, v := range hotSearchTypeMap {
		result = append(result, map[string]string{
			"type": k,
			"name": v,
		})
	}
	return result
}
