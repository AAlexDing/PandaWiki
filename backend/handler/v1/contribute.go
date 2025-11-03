package v1

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/consts"
	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/middleware"
	"github.com/chaitin/panda-wiki/repo/pg"
)

type ContributeHandler struct {
	*handler.BaseHandler
	contributeRepo *pg.ContributeRepo
	logger         *log.Logger
}

func NewContributeHandler(e *echo.Echo, baseHandler *handler.BaseHandler, contributeRepo *pg.ContributeRepo, logger *log.Logger, auth middleware.AuthMiddleware) *ContributeHandler {
	h := &ContributeHandler{
		BaseHandler:    baseHandler,
		contributeRepo: contributeRepo,
		logger:         logger.WithModule("handler.v1.contribute"),
	}

	// 注册路由
	e.GET("/api/pro/v1/contribute/list", h.GetContributeList, auth.Authorize)
	e.GET("/api/pro/v1/contribute/detail", h.GetContributeDetail, auth.Authorize)
	e.POST("/api/pro/v1/contribute/audit", h.AuditContribute, auth.Authorize)
	e.POST("/api/pro/v1/contribute/approve", h.ApproveContribute, auth.Authorize)
	e.POST("/api/pro/v1/contribute/reject", h.RejectContribute, auth.Authorize)
	e.DELETE("/api/pro/v1/contribute/delete", h.DeleteContribute, auth.Authorize)

	return h
}

// GetContributeList 获取贡献列表
//
//	@Summary		Get contribute list
//	@Description	Get contribute list for knowledge base
//	@Tags			contribute
//	@Accept			json
//	@Produce		json
//	@Param			kb_id		query		string	true	"Knowledge Base ID"
//	@Param			page		query		int		false	"Page"
//	@Param			per_page	query		int		false	"Per page"
//	@Success		200			{object}	domain.PWResponse{data=GetContributeListResp}
//	@Router			/api/pro/v1/contribute/list [get]
//	@Security		bearerAuth
func (h *ContributeHandler) GetContributeList(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	page := 1
	perPage := 10
	if p := c.QueryParam("page"); p != "" {
		if err := echo.QueryParamsBinder(c).Int("page", &page).BindError(); err != nil {
			page = 1
		}
	}
	if pp := c.QueryParam("per_page"); pp != "" {
		if err := echo.QueryParamsBinder(c).Int("per_page", &perPage).BindError(); err != nil {
			perPage = 10
		}
	}

	contributes, total, err := h.contributeRepo.GetListByKBID(c.Request().Context(), kbID, page, perPage)
	if err != nil {
		h.logger.Error("get contribute list failed", log.Error(err))
		return h.NewResponseWithError(c, "get contribute list failed", err)
	}

	// 转换为响应格式，添加地理位置信息和用户信息
	items := make([]*ContributeItem, len(contributes))
	for i, contrib := range contributes {
		// node_name 逻辑：新增时用用户提交的name，编辑时用JOIN的node_name
		nodeName := contrib.NodeName
		if contrib.Type == consts.ContributeTypeAdd {
			nodeName = contrib.Name // 新增时使用用户提交的标题
		}

		// 提取用户信息
		authName := ""
		avatar := ""
		if contrib.UserInfo != nil {
			authName = contrib.UserInfo.Username
			avatar = contrib.UserInfo.AvatarUrl
		}

		items[i] = &ContributeItem{
			Id:             contrib.Id,
			AuthId:         contrib.AuthId,
			AuthName:       authName,
			Avatar:         avatar,
			KBId:           contrib.KBId,
			Status:         contrib.Status,
			Type:           contrib.Type,
			NodeId:         contrib.NodeId,
			NodeName:       nodeName,
			ContributeName: contrib.Name, // 用户提交的标题
			Content:        contrib.Content,
			Meta:           contrib.Meta,
			Reason:         contrib.Reason,
			RemoteIP:       contrib.RemoteIP,
			AuditUserID:    contrib.AuditUserID,
			AuditTime:      contrib.AuditTime,
			IPAddress: &domain.IPAddress{
				IP:       contrib.RemoteIP,
				Country:  "Unknown",
				Province: "Unknown",
				City:     "Unknown",
			},
			CreatedAt: contrib.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: contrib.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return h.NewResponseWithData(c, GetContributeListResp{
		List:  items,
		Total: total,
	})
}

// GetContributeDetail 获取贡献详情
//
//	@Summary		Get contribute detail
//	@Description	Get contribute detail by ID
//	@Tags			contribute
//	@Accept			json
//	@Produce		json
//	@Param			id		query		string	true	"Contribute ID"
//	@Param			kb_id	query		string	true	"Knowledge Base ID"
//	@Success		200		{object}	domain.PWResponse{data=ContributeItem}
//	@Router			/api/pro/v1/contribute/detail [get]
//	@Security		bearerAuth
func (h *ContributeHandler) GetContributeDetail(c echo.Context) error {
	id := c.QueryParam("id")
	kbID := c.QueryParam("kb_id")

	if id == "" || kbID == "" {
		return h.NewResponseWithError(c, "id and kb_id are required", nil)
	}

	contrib, err := h.contributeRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("get contribute detail failed", log.Error(err))
		return h.NewResponseWithError(c, "get contribute detail failed", err)
	}

	// 验证 kb_id 是否匹配
	if contrib.KBId != kbID {
		return h.NewResponseWithError(c, "contribute not found", nil)
	}

	// node_name 逻辑：新增时用用户提交的name，编辑时用JOIN的node_name
	nodeName := contrib.NodeName
	if contrib.Type == consts.ContributeTypeAdd {
		nodeName = contrib.Name // 新增时使用用户提交的标题
	}

	// 提取用户信息
	authName := ""
	avatar := ""
	if contrib.UserInfo != nil {
		authName = contrib.UserInfo.Username
		avatar = contrib.UserInfo.AvatarUrl
	}

	// 转换为响应格式
	item := &ContributeItem{
		Id:             contrib.Id,
		AuthId:         contrib.AuthId,
		AuthName:       authName,
		Avatar:         avatar,
		KBId:           contrib.KBId,
		Status:         contrib.Status,
		Type:           contrib.Type,
		NodeId:         contrib.NodeId,
		NodeName:       nodeName,
		ContributeName: contrib.Name,
		Content:        contrib.Content,
		Meta:           contrib.Meta,
		Reason:         contrib.Reason,
		RemoteIP:       contrib.RemoteIP,
		AuditUserID:    contrib.AuditUserID,
		AuditTime:      contrib.AuditTime,
		IPAddress: &domain.IPAddress{
			IP:       contrib.RemoteIP,
			Country:  "Unknown",
			Province: "Unknown",
			City:     "Unknown",
		},
		CreatedAt: contrib.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: contrib.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return h.NewResponseWithData(c, item)
}

// AuditContribute 审核贡献（统一接口）
//
//	@Summary		Audit contribute
//	@Description	Audit a contribute (approve or reject)
//	@Tags			contribute
//	@Accept			json
//	@Produce		json
//	@Param			body	body		AuditContributeReq	true	"Audit contribute request"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/api/pro/v1/contribute/audit [post]
//	@Security		bearerAuth
func (h *ContributeHandler) AuditContribute(c echo.Context) error {
	var req AuditContributeReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}

	if req.ID == "" || req.KBID == "" || req.Status == "" {
		return h.NewResponseWithError(c, "id, kb_id and status are required", nil)
	}

	// 验证状态
	if req.Status != consts.ContributeStatusApproved && req.Status != consts.ContributeStatusRejected {
		return h.NewResponseWithError(c, "invalid status", nil)
	}

	// 获取当前用户 ID
	authInfo := domain.GetAuthInfoFromCtx(c.Request().Context())
	if authInfo == nil {
		return h.NewResponseWithError(c, "unauthorized", nil)
	}

	reason := req.Reason
	if reason == "" {
		if req.Status == consts.ContributeStatusApproved {
			reason = "approved"
		} else {
			reason = "rejected"
		}
	}

	if err := h.contributeRepo.UpdateStatus(c.Request().Context(), req.ID, req.Status, authInfo.UserId, reason); err != nil {
		h.logger.Error("audit contribute failed", log.Error(err))
		return h.NewResponseWithError(c, "audit contribute failed", err)
	}

	message := "contribute approved successfully"
	if req.Status == consts.ContributeStatusRejected {
		message = "contribute rejected successfully"
	}

	return h.NewResponseWithData(c, map[string]string{"message": message})
}

// ApproveContribute 批准贡献
//
//	@Summary		Approve contribute
//	@Description	Approve a contribute
//	@Tags			contribute
//	@Accept			json
//	@Produce		json
//	@Param			body	body		UpdateContributeReq	true	"Approve contribute request"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/api/pro/v1/contribute/approve [post]
//	@Security		bearerAuth
func (h *ContributeHandler) ApproveContribute(c echo.Context) error {
	var req UpdateContributeReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}

	if req.ID == "" {
		return h.NewResponseWithError(c, "id is required", nil)
	}

	// 获取当前用户 ID
	authInfo := domain.GetAuthInfoFromCtx(c.Request().Context())
	if authInfo == nil {
		return h.NewResponseWithError(c, "unauthorized", nil)
	}

	reason := req.Reason
	if reason == "" {
		reason = "approved"
	}

	if err := h.contributeRepo.UpdateStatus(c.Request().Context(), req.ID, consts.ContributeStatusApproved, authInfo.UserId, reason); err != nil {
		h.logger.Error("approve contribute failed", log.Error(err))
		return h.NewResponseWithError(c, "approve contribute failed", err)
	}

	return h.NewResponseWithData(c, map[string]string{"message": "contribute approved successfully"})
}

// RejectContribute 拒绝贡献
//
//	@Summary		Reject contribute
//	@Description	Reject a contribute
//	@Tags			contribute
//	@Accept			json
//	@Produce		json
//	@Param			body	body		UpdateContributeReq	true	"Reject contribute request"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/api/pro/v1/contribute/reject [post]
//	@Security		bearerAuth
func (h *ContributeHandler) RejectContribute(c echo.Context) error {
	var req UpdateContributeReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}

	if req.ID == "" || req.Reason == "" {
		return h.NewResponseWithError(c, "id and reason are required", nil)
	}

	// 获取当前用户 ID
	authInfo := domain.GetAuthInfoFromCtx(c.Request().Context())
	if authInfo == nil {
		return h.NewResponseWithError(c, "unauthorized", nil)
	}

	if err := h.contributeRepo.UpdateStatus(c.Request().Context(), req.ID, consts.ContributeStatusRejected, authInfo.UserId, req.Reason); err != nil {
		h.logger.Error("reject contribute failed", log.Error(err))
		return h.NewResponseWithError(c, "reject contribute failed", err)
	}

	return h.NewResponseWithData(c, map[string]string{"message": "contribute rejected successfully"})
}

// DeleteContribute 删除贡献
//
//	@Summary		Delete contribute
//	@Description	Delete a contribute
//	@Tags			contribute
//	@Accept			json
//	@Produce		json
//	@Param			id		query		string	true	"Contribute ID"
//	@Param			kb_id	query		string	true	"Knowledge Base ID"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/api/pro/v1/contribute/delete [delete]
//	@Security		bearerAuth
func (h *ContributeHandler) DeleteContribute(c echo.Context) error {
	id := c.QueryParam("id")
	kbID := c.QueryParam("kb_id")

	if id == "" || kbID == "" {
		return h.NewResponseWithError(c, "id and kb_id are required", nil)
	}

	if err := h.contributeRepo.Delete(c.Request().Context(), id, kbID); err != nil {
		h.logger.Error("delete contribute failed", log.Error(err))
		return h.NewResponseWithError(c, "delete contribute failed", err)
	}

	return h.NewResponseWithData(c, map[string]string{"message": "contribute deleted successfully"})
}

// GetContributeListResp 获取贡献列表响应
type GetContributeListResp struct {
	List  []*ContributeItem `json:"list"`
	Total int64             `json:"total"`
}

// ContributeItem 贡献列表项
type ContributeItem struct {
	Id             string                  `json:"id"`
	AuthId         *int64                  `json:"auth_id"`
	AuthName       string                  `json:"auth_name,omitempty"` // 用户名称
	Avatar         string                  `json:"avatar,omitempty"`    // 用户头像
	KBId           string                  `json:"kb_id"`
	Status         consts.ContributeStatus `json:"status"`
	Type           consts.ContributeType   `json:"type"`
	NodeId         string                  `json:"node_id"`
	NodeName       string                  `json:"node_name"`       // 文档标题（编辑时为原文档名，新增时为用户提交的name）
	ContributeName string                  `json:"contribute_name"` // 用户提交的标题（前端使用的字段名）
	Content        string                  `json:"content"`
	Meta           domain.NodeMeta         `json:"meta"`
	Reason         string                  `json:"reason"`    // 提交说明/审核原因
	RemoteIP       string                  `json:"remote_ip"` // 远程IP
	AuditUserID    string                  `json:"audit_user_id"`
	AuditTime      *time.Time              `json:"audit_time"`
	IPAddress      *domain.IPAddress       `json:"ip_address"`
	CreatedAt      string                  `json:"created_at"`
	UpdatedAt      string                  `json:"updated_at"`
}

// AuditContributeReq 审核贡献请求
type AuditContributeReq struct {
	ID     string                  `json:"id" validate:"required"`
	KBID   string                  `json:"kb_id" validate:"required"`
	Status consts.ContributeStatus `json:"status" validate:"required"` // approved 或 rejected
	Reason string                  `json:"reason"`
}

// UpdateContributeReq 更新贡献状态请求
type UpdateContributeReq struct {
	ID     string `json:"id" validate:"required"`
	Reason string `json:"reason"`
}
