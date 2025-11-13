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

	// 系统状态使用数据运营权限
	group := echo.Group("/api/v1/system", h.auth.Authorize, auth.ValidateKBUserPerm(consts.UserKBPermissionDataOperate))

	group.GET("", h.GetSystem)
	group.GET("/logs/:containerName", h.GetContainerLogs)
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

// GetContainerLogs 获取容器日志
//
//	@Summary		获取容器日志
//	@Description	获取指定容器的分页日志
//	@Tags			system
//	@Accept			json
//	@Produce		json
//	@Security		bearerAuth
//	@Param			containerName	path	string	true	"容器名称"
//	@Param			page			query	int		false	"页码，从1开始"	default(1)
//	@Param			limit			query	int		false	"每页大小"	default(100)
//	@Success		200				{object}	domain.PWResponse{data=v1.ContainerLogsResp}
//	@Router			/api/v1/system/logs/{containerName} [get]
func (h *SystemHandler) GetContainerLogs(c echo.Context) error {
	var req v1.ContainerLogsReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request parameters", err)
	}

	// 从路径参数获取容器名称
	req.ContainerName = c.Param("containerName")
	if req.ContainerName == "" {
		return h.NewResponseWithError(c, "container name is required", nil)
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 100
	}

	logs, err := h.usecase.GetContainerLogs(c.Request().Context(), req.ContainerName, req.Page, req.Limit)
	if err != nil {
		h.logger.Error("get container logs failed",
			log.String("container", req.ContainerName),
			log.Int("page", req.Page),
			log.Int("limit", req.Limit),
			log.Error(err))
		return h.NewResponseWithError(c, "get container logs failed", err)
	}

	return h.NewResponseWithData(c, logs)
}
