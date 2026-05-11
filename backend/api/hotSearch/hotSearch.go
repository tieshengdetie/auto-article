package hotSearch

import (
	"AutoArticle/api/hotSearch/dto"
	"AutoArticle/service"
	"AutoArticle/utils"

	"github.com/gin-gonic/gin"
)

type HotSearch struct {
	hotSearchService *service.HotSearchServer
}

func NewHotSearch() *HotSearch {
	return &HotSearch{
		hotSearchService: &service.HotSearchServer{},
	}
}

// HandleGetHotSearch 获取热搜数据
//
//	@Description: 根据热搜类型获取对应的热搜榜单
//	@receiver h
//	@param c
func (h *HotSearch) HandleGetHotSearch(c *gin.Context) {
	var req dto.HotSearchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Err(c, 400, err.Error())
		return
	}

	resp, err := h.hotSearchService.GetHotSearch(c, req)
	if err != nil {
		utils.Err(c, 500, err.Error())
		return
	}

	utils.Success(c, resp)
}

// HandleGetHotSearchTypes 获取热搜类型列表
//
//	@Description: 获取系统支持的所有热搜类型
//	@receiver h
//	@param c
func (h *HotSearch) HandleGetHotSearchTypes(c *gin.Context) {
	types := h.hotSearchService.GetHotSearchTypes()
	utils.Success(c, types)
}
