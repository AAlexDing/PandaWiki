package pro

import (
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/middleware"
)

type DocumentFeedbackHandler struct {
	*handler.BaseHandler
	logger *log.Logger
}

func NewDocumentFeedbackHandler(e *echo.Echo, baseHandler *handler.BaseHandler, logger *log.Logger, auth middleware.AuthMiddleware) *DocumentFeedbackHandler {
	h := &DocumentFeedbackHandler{
		BaseHandler: baseHandler,
		logger:      logger.WithModule("handler.pro.document_feedback"),
	}

	// 注册路由
	e.POST("/api/pro/v1/document/feedback", h.SubmitFeedback, auth.Authorize)
	e.GET("/api/pro/v1/document/list", h.GetDocumentList, auth.Authorize)

	return h
}

// SubmitFeedback 提交文档反馈
//
//	@Summary		提交文档反馈
//	@Description	用户提交文档反馈意见
//	@Tags			Document
//	@Accept			json
//	@Produce		json
//	@Param			body	body		SubmitFeedbackReq	true	"Feedback request"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/api/pro/v1/document/feedback [post]
//	@Security		bearerAuth
func (h *DocumentFeedbackHandler) SubmitFeedback(c echo.Context) error {
	var req SubmitFeedbackReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}

	if req.NodeID == "" || req.Content == "" {
		return h.NewResponseWithError(c, "node_id and content are required", nil)
	}

	// TODO: 保存反馈到数据库
	h.logger.Info("Document feedback", "node_id", req.NodeID, "content_length", len(req.Content))

	return h.NewResponseWithData(c, map[string]string{"message": "Feedback submitted successfully"})
}

// GetDocumentList 获取文档反馈列表
//
//	@Summary		获取文档反馈列表
//	@Description	获取知识库的文档反馈列表
//	@Tags			Document
//	@Accept			json
//	@Produce		json
//	@Param			kb_id		query		string	true	"Knowledge Base ID"
//	@Param			page		query		int		false	"Page"
//	@Param			per_page	query		int		false	"Per page"
//	@Success		200			{object}	domain.PWResponse{data=GetDocumentListResp}
//	@Router			/api/pro/v1/document/list [get]
//	@Security		bearerAuth
func (h *DocumentFeedbackHandler) GetDocumentList(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	// TODO: 从数据库查询文档反馈列表
	return h.NewResponseWithData(c, GetDocumentListResp{
		Data:  []DocumentFeedbackItem{},
		Total: 0,
	})
}

// SubmitFeedbackReq 提交反馈请求
type SubmitFeedbackReq struct {
	NodeID               string `json:"node_id" validate:"required"`
	Content              string `json:"content" validate:"required"`
	CorrectionSuggestion string `json:"correction_suggestion"` // 修正建议
	Email                string `json:"email"`
}

// GetDocumentListResp 文档反馈列表响应
type GetDocumentListResp struct {
	Data  []DocumentFeedbackItem `json:"data"`
	Total int64                  `json:"total"`
}

// DocumentFeedbackItem 文档反馈项
type DocumentFeedbackItem struct {
	ID                   string                   `json:"id"`
	NodeID               string                   `json:"node_id"`
	NodeName             string                   `json:"node_name"`
	KBId                 string                   `json:"kb_id"`
	UserID               string                   `json:"user_id"`
	Content              string                   `json:"content"`
	CorrectionSuggestion string                   `json:"correction_suggestion"` // 修正建议
	Info                 *DocumentFeedbackInfo    `json:"info"`                  // 用户信息
	IPAddress            *domain.IPAddress        `json:"ip_address"`
	CreatedAt            string                   `json:"created_at"`
}

// DocumentFeedbackInfo 反馈用户信息
type DocumentFeedbackInfo struct {
	AuthUserID int64  `json:"auth_user_id"` // 用户 ID
	UserName   string `json:"user_name"`    // 用户名
	Avatar     string `json:"avatar"`       // 头像
	Email      string `json:"email"`        // 邮箱
	RemoteIP   string `json:"remote_ip"`    // IP 地址
	ScreenShot string `json:"screen_shot"`  // 截图
}
