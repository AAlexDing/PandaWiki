package v1

import (
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/middleware"
	"github.com/chaitin/panda-wiki/repo/pg"
)

type BlockWordHandler struct {
	*handler.BaseHandler
	blockWordRepo *pg.BlockWordRepo
	logger        *log.Logger
}

func NewBlockWordHandler(e *echo.Echo, baseHandler *handler.BaseHandler, blockWordRepo *pg.BlockWordRepo, logger *log.Logger, auth middleware.AuthMiddleware) *BlockWordHandler {
	h := &BlockWordHandler{
		BaseHandler:   baseHandler,
		blockWordRepo: blockWordRepo,
		logger:        logger.WithModule("handler.v1.block_word"),
	}

	// 注册路由
	e.GET("/api/pro/v1/block", h.GetBlockWords, auth.Authorize)
	e.POST("/api/pro/v1/block", h.UpdateBlockWords, auth.Authorize)

	return h
}

// GetBlockWords 获取屏蔽词列表
//
//	@Summary		Get block words
//	@Description	Get block words for knowledge base
//	@Tags			block
//	@Accept			json
//	@Produce		json
//	@Param			kb_id	query		string	true	"Knowledge Base ID"
//	@Success		200		{object}	domain.PWResponse{data=GetBlockWordsResp}
//	@Router			/api/pro/v1/block [get]
//	@Security		bearerAuth
func (h *BlockWordHandler) GetBlockWords(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	words, err := h.blockWordRepo.GetBlockWords(c.Request().Context(), kbID)
	if err != nil {
		h.logger.Error("get block words failed", log.Error(err))
		return h.NewResponseWithError(c, "get block words failed", err)
	}

	if words == nil {
		words = []string{}
	}

	return h.NewResponseWithData(c, GetBlockWordsResp{Words: words})
}

// UpdateBlockWords 更新屏蔽词列表
//
//	@Summary		Update block words
//	@Description	Update block words for knowledge base
//	@Tags			block
//	@Accept			json
//	@Produce		json
//	@Param			body	body		UpdateBlockWordsReq	true	"Update block words request"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/api/pro/v1/block [post]
//	@Security		bearerAuth
func (h *BlockWordHandler) UpdateBlockWords(c echo.Context) error {
	var req UpdateBlockWordsReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}

	if req.KBID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	if req.Words == nil {
		req.Words = []string{}
	}

	if err := h.blockWordRepo.UpdateBlockWords(c.Request().Context(), req.KBID, req.Words); err != nil {
		h.logger.Error("update block words failed", log.Error(err))
		return h.NewResponseWithError(c, "update block words failed", err)
	}

	return h.NewResponseWithData(c, map[string]string{"message": "block words updated successfully"})
}

// UpdateBlockWordsReq 更新屏蔽词请求
type UpdateBlockWordsReq struct {
	KBID  string   `json:"kb_id" validate:"required"`
	Words []string `json:"words"`
}

// GetBlockWordsResp 获取屏蔽词响应
type GetBlockWordsResp struct {
	Words []string `json:"words"`
}
