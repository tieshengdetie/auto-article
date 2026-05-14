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

// HandleGetHotSearch 获取指定类型热搜
func (h *HotSearch) HandleGetHotSearch(c *gin.Context) {
	var req dto.HotSearchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Err(c, 400, err.Error())
		return
	}

	resp, err := h.hotSearchService.GetHotSearch(c.Request.Context(), req)
	if err != nil {
		utils.Err(c, 500, err.Error())
		return
	}

	utils.Success(c, resp)
}

// HandleGetHotSearchTypes 获取支持的热搜类型
func (h *HotSearch) HandleGetHotSearchTypes(c *gin.Context) {
	utils.Success(c, h.hotSearchService.GetHotSearchTypes())
}
