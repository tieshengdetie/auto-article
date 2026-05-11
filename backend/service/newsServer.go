package service

import (
	"AutoArticle/api/news/dto"
	"AutoArticle/api/news/vo"
	"AutoArticle/common/sendHttp"
	"context"
)

type NewsServer struct {
}

// 新闻频道映射
var newsChannelMap = []vo.NewsChannel{
	{ChannelID: 46, ChannelName: "宠物新闻", Description: "宠物及相关产业新闻资讯"},
	{ChannelID: 45, ChannelName: "电竞资讯", Description: "电子竞技新闻资讯"},
	{ChannelID: 43, ChannelName: "女性新闻", Description: "新浪女性新闻频道"},
	{ChannelID: 42, ChannelName: "垃圾分类新闻", Description: "垃圾分类新闻资讯"},
	{ChannelID: 41, ChannelName: "环保资讯", Description: "人民网环保新闻资讯"},
	{ChannelID: 40, ChannelName: "影视资讯", Description: "影娱乐资讯新闻"},
	{ChannelID: 36, ChannelName: "科学探索", Description: "探索宇宙和科学的真相"},
	{ChannelID: 35, ChannelName: "汽车新闻", Description: "汽车行业新闻资讯"},
	{ChannelID: 34, ChannelName: "互联网资讯", Description: "互联网行业资讯新闻"},
	{ChannelID: 33, ChannelName: "动漫资讯", Description: "全网热点动漫资讯，带你了解二次元世界"},
	{ChannelID: 32, ChannelName: "财经新闻", Description: "财经资讯，了解身边的经济大事"},
	{ChannelID: 31, ChannelName: "游戏资讯", Description: "网络游戏相关每日精选新闻"},
	{ChannelID: 30, ChannelName: "CBA新闻", Description: "中国男子职业篮球赛资讯等"},
	{ChannelID: 29, ChannelName: "人工智能", Description: "AI人工智能行业相关新闻资讯"},
	{ChannelID: 27, ChannelName: "军事新闻", Description: "军事资讯、军情动态、科技发展等"},
	{ChannelID: 26, ChannelName: "足球新闻", Description: "国足资讯、国足明星动态等"},
	{ChannelID: 22, ChannelName: "IT资讯", Description: "IT行业相关新闻资讯"},
	{ChannelID: 21, ChannelName: "VR科技", Description: "VR虚拟现实相关新闻资讯"},
	{ChannelID: 20, ChannelName: "NBA新闻", Description: "NBA新闻动态、篮球赛等"},
	{ChannelID: 18, ChannelName: "旅游资讯", Description: "旅游、周边、景点新闻资讯"},
	{ChannelID: 17, ChannelName: "健康知识", Description: "健康知识、养生、中西医资讯"},
	{ChannelID: 13, ChannelName: "科技新闻", Description: "信息科技行业新闻、物理科技"},
	{ChannelID: 12, ChannelName: "体育新闻", Description: "国内外体育、体育明星动态等"},
	{ChannelID: 10, ChannelName: "娱乐新闻", Description: "明星花边、探班、娱乐活动等"},
	{ChannelID: 8, ChannelName: "国际新闻", Description: "返回国际新闻资讯"},
	{ChannelID: 7, ChannelName: "国内新闻", Description: "返回国内新闻资讯"},
	{ChannelID: 5, ChannelName: "社会新闻", Description: "返回社会新闻资讯"},
}

// GetAllNews 获取指定分类新闻
//
//	@Description: 获取指定分类的新闻列表，支持关键词搜索和翻页
//	@receiver n
//	@param ctx
//	@param req
//	@return resp
//	@return err
func (n *NewsServer) GetAllNews(ctx context.Context, req dto.AllNewsReq) (resp *vo.NewsResp, err error) {
	// 设置默认值
	if req.Num <= 0 {
		req.Num = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	data, err := sendHttp.TianApi.GetAllNews(ctx, req.Col, req.Num, req.Page, req.Rand, req.Word)
	if err != nil {
		return nil, err
	}

	resp = &vo.NewsResp{
		Allnum:  data.Result.Allnum,
		Curpage: data.Result.Curpage,
		List:    make([]vo.NewsItem, 0, len(data.Result.List)),
	}

	for _, item := range data.Result.List {
		resp.List = append(resp.List, vo.NewsItem{
			ID:          item.ID,
			Url:         item.Url,
			Ctime:       item.Ctime,
			Title:       item.Title,
			PicUrl:      item.PicUrl,
			Source:      item.Source,
			Description: item.Description,
		})
	}

	return resp, nil
}

// GetGeneralNews 获取综合新闻
//
//	@Description: 获取综合新闻列表，支持关键词搜索、翻页和指定来源
//	@receiver n
//	@param ctx
//	@param req
//	@return resp
//	@return err
func (n *NewsServer) GetGeneralNews(ctx context.Context, req dto.GeneralNewsReq) (resp *vo.NewsResp, err error) {
	// 设置默认值
	if req.Num <= 0 {
		req.Num = 10
	}

	data, err := sendHttp.TianApi.GetGeneralNews(ctx, req.Num, req.Page, req.Rand, req.Word, req.Source)
	if err != nil {
		return nil, err
	}

	resp = &vo.NewsResp{
		Allnum:  data.Result.Allnum,
		Curpage: data.Result.Curpage,
		List:    make([]vo.NewsItem, 0, len(data.Result.List)),
	}

	for _, item := range data.Result.List {
		resp.List = append(resp.List, vo.NewsItem{
			ID:          item.ID,
			Url:         item.Url,
			Ctime:       item.Ctime,
			Title:       item.Title,
			PicUrl:      item.PicUrl,
			Source:      item.Source,
			Description: item.Description,
		})
	}

	return resp, nil
}

// GetNewsChannels 获取新闻频道列表
//
//	@Description: 获取所有支持的新闻频道列表
//	@receiver n
//	@return channels
func (n *NewsServer) GetNewsChannels() []vo.NewsChannel {
	return newsChannelMap
}
