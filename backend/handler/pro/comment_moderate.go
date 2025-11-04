package pro

import (
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/middleware"
)

type CommentModerateHandler struct {
	*handler.BaseHandler
	logger *log.Logger
}

func NewCommentModerateHandler(e *echo.Echo, baseHandler *handler.BaseHandler, logger *log.Logger, auth middleware.AuthMiddleware) *CommentModerateHandler {
	h := &CommentModerateHandler{
		BaseHandler: baseHandler,
		logger:      logger.WithModule("handler.pro.comment_moderate"),
	}

	// 注册路由
	e.POST("/api/pro/v1/comment_moderate", h.ModerateComment, auth.Authorize)

	return h
}

// ModerateComment 审核评论
//
//	@Summary		审核评论
//	@Description	审核通过或拒绝评论
//	@Tags			Comment
//	@Accept			json
//	@Produce		json
//	@Param			body	body		ModerateCommentReq	true	"Moderate request"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/api/pro/v1/comment_moderate [post]
//	@Security		bearerAuth
func (h *CommentModerateHandler) ModerateComment(c echo.Context) error {
	var req ModerateCommentReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}

	if req.CommentID == "" || req.Action == "" {
		return h.NewResponseWithError(c, "comment_id and action are required", nil)
	}

	if req.Action != "approve" && req.Action != "reject" {
		return h.NewResponseWithError(c, "action must be approve or reject", nil)
	}

	// TODO: 更新评论状态到数据库
	h.logger.Info("Moderate comment", "comment_id", req.CommentID, "action", req.Action)

	return h.NewResponseWithData(c, map[string]string{"message": "Comment moderated successfully"})
}

// ModerateCommentReq 审核评论请求
type ModerateCommentReq struct {
	CommentID string `json:"comment_id" validate:"required"`
	Action    string `json:"action" validate:"required"` // approve, reject
	Reason    string `json:"reason"`                     // 拒绝原因（reject 时必填）
}
