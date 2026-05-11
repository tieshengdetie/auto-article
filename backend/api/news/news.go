package news

import (
	"AutoArticle/api/news/dto"
	"AutoArticle/api/news/vo"
	"AutoArticle/service"
	"AutoArticle/utils"
	"context"

	"github.com/gin-gonic/gin"
)

type News struct {
	newsService *service.NewsServer
}

func NewNews() *News {
	return &News{
		newsService: &service.NewsServer{},
	}
}

// GetAllNews 获取指定分类新闻
//
//	@Description: 获取指定分类的新闻列表，支持关键词搜索和翻页
//	@receiver n
//	@param ctx
//	@param req
//	@return resp
//	@return err
func (n *News) GetAllNews(ctx context.Context, req dto.AllNewsReq) (resp *vo.NewsResp, err error) {
	return n.newsService.GetAllNews(ctx, req)
}

// GetGeneralNews 获取综合新闻
//
//	@Description: 获取综合新闻列表，支持关键词搜索、翻页和指定来源
//	@receiver n
//	@param ctx
//	@param req
//	@return resp
//	@return err
func (n *News) GetGeneralNews(ctx context.Context, req dto.GeneralNewsReq) (resp *vo.NewsResp, err error) {
	return n.newsService.GetGeneralNews(ctx, req)
}

// GetNewsChannels 获取新闻频道列表
//
//	@Description: 获取所有支持的新闻频道列表
//	@receiver n
//	@return channels
func (n *News) GetNewsChannels() []vo.NewsChannel {
	return n.newsService.GetNewsChannels()
}

// HandleGetAllNews 获取指定分类新闻 Handler
func (n *News) HandleGetAllNews(c *gin.Context) {
	var req dto.AllNewsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Err(c, 400, err.Error())
		return
	}

	resp, err := n.GetAllNews(c.Request.Context(), req)
	if err != nil {
		utils.Err(c, 500, err.Error())
		return
	}

	utils.Success(c, resp)
}

// HandleGetGeneralNews 获取综合新闻 Handler
func (n *News) HandleGetGeneralNews(c *gin.Context) {
	var req dto.GeneralNewsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Err(c, 400, err.Error())
		return
	}

	resp, err := n.GetGeneralNews(c.Request.Context(), req)
	if err != nil {
		utils.Err(c, 500, err.Error())
		return
	}

	utils.Success(c, resp)
}

// HandleGetNewsChannels 获取新闻频道列表 Handler
func (n *News) HandleGetNewsChannels(c *gin.Context) {
	channels := n.GetNewsChannels()
	utils.Success(c, channels)
}
