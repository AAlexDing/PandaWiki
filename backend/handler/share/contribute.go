package share

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/consts"
	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/repo/pg"
)

type ShareContributeHandler struct {
	*handler.BaseHandler
	contributeRepo *pg.ContributeRepo
	logger         *log.Logger
}

func NewShareContributeHandler(e *echo.Echo, baseHandler *handler.BaseHandler, contributeRepo *pg.ContributeRepo, logger *log.Logger) *ShareContributeHandler {
	h := &ShareContributeHandler{
		BaseHandler:    baseHandler,
		contributeRepo: contributeRepo,
		logger:         logger.WithModule("handler.share.contribute"),
	}

	// 注册路由
	e.POST("/share/pro/v1/contribute/submit", h.SubmitContribute)

	return h
}

// SubmitContribute 提交贡献
//
//	@Summary		Submit contribute
//	@Description	Submit a new contribute for knowledge base
//	@Tags			contribute
//	@Accept			json
//	@Produce		json
//	@Param			body	body		SubmitContributeReq	true	"Submit contribute request"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/share/pro/v1/contribute/submit [post]
func (h *ShareContributeHandler) SubmitContribute(c echo.Context) error {
	var req SubmitContributeReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}

	// 从请求头获取 kb_id（分享页面统一使用 X-KB-ID 头）
	if req.KBID == "" {
		req.KBID = c.Request().Header.Get("X-KB-ID")
	}

	// 尝试从 query 参数获取 kb_id（如果请求体和头都没有）
	if req.KBID == "" {
		req.KBID = c.QueryParam("kb_id")
	}

	// 基本验证
	if req.Content == "" || req.Type == "" {
		return h.NewResponseWithError(c, "content and type are required", nil)
	}

	if req.KBID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	// 获取远程 IP
	remoteIP := c.RealIP()

	// 获取认证信息（可能为空，因为分享页面可能允许匿名贡献）
	var authID *int64
	authInfo := domain.GetAuthInfoFromCtx(c.Request().Context())
	if authInfo != nil && authInfo.UserId != "" {
		// 这里需要根据 UserId 查询 auth_id，暂时留空
		// 如果需要可以添加 AuthRepo 查询
	}

	contribute := &domain.Contribute{
		Id:          uuid.New().String(),
		AuthId:      authID,
		KBId:        req.KBID,
		Status:      consts.ContributeStatusPending,
		Type:        req.Type,
		NodeId:      req.NodeID,
		Name:        req.Name,
		Content:     req.Content,
		Meta:        req.Meta,
		Reason:      req.Reason, // 用户提交时的说明
		AuditUserID: "",
		RemoteIP:    remoteIP,
	}

	if err := h.contributeRepo.Create(c.Request().Context(), contribute); err != nil {
		h.logger.Error("create contribute failed", log.Error(err))
		return h.NewResponseWithError(c, "create contribute failed", err)
	}

	return h.NewResponseWithData(c, map[string]string{
		"id":      contribute.Id,
		"message": "contribute submitted successfully",
	})
}

// SubmitContributeReq 提交贡献请求
type SubmitContributeReq struct {
	KBID    string                `json:"kb_id" validate:"required"`
	Type    consts.ContributeType `json:"type" validate:"required"` // add 或 edit
	NodeID  string                `json:"node_id"`                  // 编辑时需要
	Name    string                `json:"name"`                     // 新增时的标题
	Content string                `json:"content" validate:"required"`
	Meta    domain.NodeMeta       `json:"meta"`
	Reason  string                `json:"reason"` // 提交说明
}
