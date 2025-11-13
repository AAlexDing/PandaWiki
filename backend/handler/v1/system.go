package v1

import (
	"github.com/labstack/echo/v4"

	v1 "github.com/chaitin/panda-wiki/api/system/v1"
	"github.com/chaitin/panda-wiki/consts"
	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/middleware"
	"github.com/chaitin/panda-wiki/usecase"
)

type SystemHandler struct {
	*handler.BaseHandler
	usecase *usecase.SystemUseCase
	auth    middleware.AuthMiddleware
	logger  *log.Logger
}

func NewSystemHandler(baseHandler *handler.BaseHandler, echo *echo.Echo, usecase *usecase.SystemUseCase, logger *log.Logger, auth middleware.AuthMiddleware) *SystemHandler {
	h := &SystemHandler{
		BaseHandler: baseHandler,
		usecase:     usecase,
		auth:        auth,
		logger:      logger.WithModule("handler.v1.system"),
	}

	group := echo.Group("/api/v1/system", h.auth.Authorize, auth.ValidateKBUserPerm(consts.UserKBPermissionFullControl))

	group.GET("", h.GetSystem)
	return h
}

// GetSystem 获取系统状态
//
//	@Summary		获取系统状态
//	@Description	获取系统状态（文档、学习、系统组件）
//	@Tags			system
//	@Accept			json
//	@Produce		json
//	@Security		bearerAuth
//	@Param			kb_id	query		string	true	"知识库ID"
//	@Success		200		{object}	domain.PWResponse{data=v1.SystemResp}
//	@Router			/api/v1/system [get]
func (h *SystemHandler) GetSystem(c echo.Context) error {
	var req v1.SystemReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request parameters", err)
	}

	if err := c.Validate(&req); err != nil {
		return h.NewResponseWithError(c, "validation failed", err)
	}

	system, err := h.usecase.GetSystem(c.Request().Context(), req.KbID)
	if err != nil {
		return h.NewResponseWithError(c, "get system failed", err)
	}
	return h.NewResponseWithData(c, system)
}
