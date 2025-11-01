package v1

import (
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/middleware"
	"github.com/chaitin/panda-wiki/repo/pg"
)

type PromptHandler struct {
	*handler.BaseHandler
	promptRepo *pg.PromptRepo
	logger     *log.Logger
}

func NewPromptHandler(e *echo.Echo, baseHandler *handler.BaseHandler, promptRepo *pg.PromptRepo, logger *log.Logger, auth middleware.AuthMiddleware) *PromptHandler {
	h := &PromptHandler{
		BaseHandler: baseHandler,
		promptRepo:  promptRepo,
		logger:      logger.WithModule("handler.v1.prompt"),
	}

	// 注册路由
	group := e.Group("/api/pro/v1/prompt", auth.Authorize)
	group.GET("", h.GetPrompt)
	group.POST("", h.UpdatePrompt)

	return h
}

// GetPrompt 获取提示词
//
//	@Summary		Get prompt
//	@Description	Get system prompt for knowledge base
//	@Tags			prompt
//	@Accept			json
//	@Produce		json
//	@Param			kb_id	query		string	true	"Knowledge Base ID"
//	@Success		200		{object}	domain.PWResponse{data=GetPromptResp}
//	@Router			/api/pro/v1/prompt [get]
//	@Security		bearerAuth
func (h *PromptHandler) GetPrompt(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	content, err := h.promptRepo.GetPrompt(c.Request().Context(), kbID)
	if err != nil {
		h.logger.Error("get prompt failed", log.Error(err))
		return h.NewResponseWithError(c, "get prompt failed", err)
	}

	// 如果没有设置自定义提示词，返回默认提示词
	if content == "" {
		content = domain.SystemPrompt
	}

	return h.NewResponseWithData(c, GetPromptResp{Content: content})
}

// UpdatePrompt 更新提示词
//
//	@Summary		Update prompt
//	@Description	Update system prompt for knowledge base
//	@Tags			prompt
//	@Accept			json
//	@Produce		json
//	@Param			body	body		UpdatePromptReq	true	"Update prompt request"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/api/pro/v1/prompt [post]
//	@Security		bearerAuth
func (h *PromptHandler) UpdatePrompt(c echo.Context) error {
	var req UpdatePromptReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}

	if req.KBID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	if err := h.promptRepo.UpdatePrompt(c.Request().Context(), req.KBID, req.Content); err != nil {
		h.logger.Error("update prompt failed", log.Error(err))
		return h.NewResponseWithError(c, "update prompt failed", err)
	}

	return h.NewResponseWithData(c, map[string]string{"message": "prompt updated successfully"})
}

// UpdatePromptReq 更新提示词请求
type UpdatePromptReq struct {
	KBID    string `json:"kb_id" validate:"required"`
	Content string `json:"content"`
}

// GetPromptResp 获取提示词响应
type GetPromptResp struct {
	Content string `json:"content"`
}
