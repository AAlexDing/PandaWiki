package share

import (
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
)

// ShareAuthProHandler Mock Share Auth Pro Handler (用于消除 404 日志)
type ShareAuthProHandler struct {
	*handler.BaseHandler
	logger *log.Logger
}

// NewShareAuthProHandler 创建 Share Auth Pro Handler 并注册路由
func NewShareAuthProHandler(e *echo.Echo, baseHandler *handler.BaseHandler, logger *log.Logger) *ShareAuthProHandler {
	h := &ShareAuthProHandler{
		BaseHandler: baseHandler,
		logger:      logger.WithModule("handler.share.pro.auth"),
	}

	// 注册路由 - 返回空数据即可
	e.GET("/share/pro/v1/auth/info", h.GetAuthInfo)

	return h
}

// GetAuthInfo 获取认证信息（空实现，消除 404）
func (h *ShareAuthProHandler) GetAuthInfo(c echo.Context) error {
	// 返回空对象，前端会理解为"没有高级认证配置"
	return h.NewResponseWithData(c, map[string]interface{}{})
}
