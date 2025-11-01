package v1

import (
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/middleware"
)

// ProMockHandler 统一的商业版 Mock Handler (用于消除 404 错误)
type ProMockHandler struct {
	*handler.BaseHandler
	logger *log.Logger
}

// NewProMockHandler 创建 Pro Mock Handler 并注册所有商业版路由
func NewProMockHandler(e *echo.Echo, baseHandler *handler.BaseHandler, logger *log.Logger, auth middleware.AuthMiddleware) *ProMockHandler {
	h := &ProMockHandler{
		BaseHandler: baseHandler,
		logger:      logger.WithModule("handler.pro.mock"),
	}

	// 注册所有商业版 API 路由，统一返回空数据
	e.GET("/api/pro/v1/contribute/list", h.GetContributeList, auth.Authorize)

	return h
}

// GetContributeList 获取贡献列表（空实现）
//
//	@Summary		Get contribute list (Mock)
//	@Description	Mock API for pro contribute list
//	@Tags			pro
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page"
//	@Param			per_page	query		int		false	"Per page"
//	@Param			kb_id		query		string	true	"Knowledge Base ID"
//	@Success		200			{object}	domain.PWResponse
//	@Router			/api/pro/v1/contribute/list [get]
//	@Security		bearerAuth
func (h *ProMockHandler) GetContributeList(c echo.Context) error {
	// 返回空列表，前端会显示"暂无数据"
	return h.NewResponseWithData(c, map[string]interface{}{
		"list":  []interface{}{},
		"total": 0,
	})
}
