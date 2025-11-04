package pro

import (
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/middleware"
)

type ReleaseHandler struct {
	*handler.BaseHandler
	logger *log.Logger
}

func NewReleaseHandler(e *echo.Echo, baseHandler *handler.BaseHandler, logger *log.Logger, auth middleware.AuthMiddleware) *ReleaseHandler {
	h := &ReleaseHandler{
		BaseHandler: baseHandler,
		logger:      logger.WithModule("handler.pro.release"),
	}

	// 注册路由
	e.GET("/api/pro/v1/node/release/list", h.GetReleaseList, auth.Authorize)
	e.GET("/api/pro/v1/node/release/detail", h.GetReleaseDetail, auth.Authorize)

	return h
}

// GetReleaseList 获取版本列表
//
//	@Summary		获取版本列表
//	@Description	获取文档的历史版本列表
//	@Tags			Release
//	@Accept			json
//	@Produce		json
//	@Param			node_id		query		string	true	"Node ID"
//	@Param			page		query		int		false	"Page"
//	@Param			per_page	query		int		false	"Per page"
//	@Success		200			{object}	domain.PWResponse{data=GetReleaseListResp}
//	@Router			/api/pro/v1/node/release/list [get]
//	@Security		bearerAuth
func (h *ReleaseHandler) GetReleaseList(c echo.Context) error {
	nodeID := c.QueryParam("node_id")
	if nodeID == "" {
		return h.NewResponseWithError(c, "node_id is required", nil)
	}

	// TODO: 从数据库查询版本列表
	return h.NewResponseWithData(c, GetReleaseListResp{
		List:  []ReleaseItem{},
		Total: 0,
	})
}

// GetReleaseDetail 获取版本详情
//
//	@Summary		获取版本详情
//	@Description	获取指定版本的详细内容
//	@Tags			Release
//	@Accept			json
//	@Produce		json
//	@Param			node_id		query		string	true	"Node ID"
//	@Param			version_id	query		string	true	"Version ID"
//	@Success		200			{object}	domain.PWResponse{data=ReleaseDetail}
//	@Router			/api/pro/v1/node/release/detail [get]
//	@Security		bearerAuth
func (h *ReleaseHandler) GetReleaseDetail(c echo.Context) error {
	nodeID := c.QueryParam("node_id")
	versionID := c.QueryParam("version_id")

	if nodeID == "" || versionID == "" {
		return h.NewResponseWithError(c, "node_id and version_id are required", nil)
	}

	// TODO: 从数据库查询版本详情
	return h.NewResponseWithData(c, ReleaseDetail{})
}

// GetReleaseListResp 版本列表响应
type GetReleaseListResp struct {
	List  []ReleaseItem `json:"list"`
	Total int64         `json:"total"`
}

// ReleaseItem 版本列表项
type ReleaseItem struct {
	ID          string `json:"id"`
	NodeID      string `json:"node_id"`
	Version     int    `json:"version"`
	Description string `json:"description"` // 版本说明
	AuthorID    int64  `json:"author_id"`
	AuthorName  string `json:"author_name"`
	CreatedAt   string `json:"created_at"`
}

// ReleaseDetail 版本详情
type ReleaseDetail struct {
	ID          string `json:"id"`
	NodeID      string `json:"node_id"`
	Version     int    `json:"version"`
	Name        string `json:"name"`
	Content     string `json:"content"`
	Description string `json:"description"`
	AuthorID    int64  `json:"author_id"`
	AuthorName  string `json:"author_name"`
	CreatedAt   string `json:"created_at"`
}
