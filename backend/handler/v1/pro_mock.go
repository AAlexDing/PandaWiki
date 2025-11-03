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

	// 所有商业版 API 已迁移到独立 Handler，这里保留作为占位符

	return h
}
