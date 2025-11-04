package pro

import (
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/middleware"
)

type AuthGroupHandler struct {
	*handler.BaseHandler
	logger *log.Logger
}

func NewAuthGroupHandler(e *echo.Echo, baseHandler *handler.BaseHandler, logger *log.Logger, auth middleware.AuthMiddleware) *AuthGroupHandler {
	h := &AuthGroupHandler{
		BaseHandler: baseHandler,
		logger:      logger.WithModule("handler.pro.auth_group"),
	}

	// 注册路由
	e.GET("/api/pro/v1/auth/group/list", h.GetList, auth.Authorize)
	e.GET("/api/pro/v1/auth/group/tree", h.GetTree, auth.Authorize)
	e.GET("/api/pro/v1/auth/group/detail", h.GetDetail, auth.Authorize)
	e.POST("/api/pro/v1/auth/group/create", h.Create, auth.Authorize)
	e.PATCH("/api/pro/v1/auth/group/update", h.Update, auth.Authorize)
	e.PATCH("/api/pro/v1/auth/group/move", h.Move, auth.Authorize)
	e.POST("/api/pro/v1/auth/group/sync", h.Sync, auth.Authorize)
	e.DELETE("/api/pro/v1/auth/group/delete", h.Delete, auth.Authorize)

	return h
}

// GetList 获取用户组列表
//
//	@Summary		获取用户组列表
//	@Description	获取指定知识库的用户组列表
//	@Tags			AuthGroup
//	@Accept			json
//	@Produce		json
//	@Param			kb_id		query		string	true	"Knowledge Base ID"
//	@Param			page		query		int		false	"Page"
//	@Param			per_page	query		int		false	"Per page"
//	@Success		200			{object}	domain.PWResponse{data=GetListResp}
//	@Router			/api/pro/v1/auth/group/list [get]
//	@Security		bearerAuth
func (h *AuthGroupHandler) GetList(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	// TODO: 从数据库查询用户组列表
	// 目前返回空列表
	return h.NewResponseWithData(c, GetListResp{
		List:  []GroupItem{},
		Total: 0,
	})
}

// GetTree 获取用户组树
//
//	@Summary		获取用户组树形结构
//	@Description	获取指定知识库的用户组树形结构
//	@Tags			AuthGroup
//	@Accept			json
//	@Produce		json
//	@Param			kb_id	query		string	true	"Knowledge Base ID"
//	@Success		200		{object}	domain.PWResponse{data=[]GroupTreeNode}
//	@Router			/api/pro/v1/auth/group/tree [get]
//	@Security		bearerAuth
func (h *AuthGroupHandler) GetTree(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	// TODO: 从数据库查询用户组树
	// 目前返回空数组
	return h.NewResponseWithData(c, []GroupTreeNode{})
}

// GetDetail 获取用户组详情
//
//	@Summary		获取用户组详情
//	@Description	获取指定用户组的详细信息
//	@Tags			AuthGroup
//	@Accept			json
//	@Produce		json
//	@Param			id		query		int		true	"Group ID"
//	@Param			kb_id	query		string	true	"Knowledge Base ID"
//	@Success		200		{object}	domain.PWResponse{data=GroupDetail}
//	@Router			/api/pro/v1/auth/group/detail [get]
//	@Security		bearerAuth
func (h *AuthGroupHandler) GetDetail(c echo.Context) error {
	id := c.QueryParam("id")
	kbID := c.QueryParam("kb_id")

	if id == "" || kbID == "" {
		return h.NewResponseWithError(c, "id and kb_id are required", nil)
	}

	// TODO: 从数据库查询用户组详情
	return h.NewResponseWithData(c, GroupDetail{})
}

// Create 创建用户组
//
//	@Summary		创建用户组
//	@Description	创建新的用户组
//	@Tags			AuthGroup
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateReq	true	"Create request"
//	@Success		200		{object}	domain.PWResponse{data=CreateResp}
//	@Router			/api/pro/v1/auth/group/create [post]
//	@Security		bearerAuth
func (h *AuthGroupHandler) Create(c echo.Context) error {
	var req CreateReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}

	if req.KBId == "" || req.Name == "" {
		return h.NewResponseWithError(c, "kb_id and name are required", nil)
	}

	// TODO: 保存到数据库
	h.logger.Info("Create group", "kb_id", req.KBId, "name", req.Name)

	return h.NewResponseWithData(c, CreateResp{
		ID:      1,
		Message: "Group created successfully",
	})
}

// Update 更新用户组
//
//	@Summary		更新用户组
//	@Description	更新用户组名称和成员
//	@Tags			AuthGroup
//	@Accept			json
//	@Produce		json
//	@Param			body	body		UpdateReq	true	"Update request"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/api/pro/v1/auth/group/update [patch]
//	@Security		bearerAuth
func (h *AuthGroupHandler) Update(c echo.Context) error {
	var req UpdateReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}

	if req.ID == 0 || req.KBId == "" {
		return h.NewResponseWithError(c, "id and kb_id are required", nil)
	}

	// TODO: 更新数据库
	h.logger.Info("Update group", "id", req.ID, "name", req.Name)

	return h.NewResponseWithData(c, map[string]string{"message": "Group updated successfully"})
}

// Move 移动用户组
//
//	@Summary		移动用户组
//	@Description	移动用户组到新的父组下
//	@Tags			AuthGroup
//	@Accept			json
//	@Produce		json
//	@Param			body	body		MoveReq	true	"Move request"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/api/pro/v1/auth/group/move [patch]
//	@Security		bearerAuth
func (h *AuthGroupHandler) Move(c echo.Context) error {
	var req MoveReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}

	if req.ID == 0 || req.KBId == "" {
		return h.NewResponseWithError(c, "id and kb_id are required", nil)
	}

	// TODO: 更新数据库
	h.logger.Info("Move group", "id", req.ID, "parent_id", req.ParentID)

	return h.NewResponseWithData(c, map[string]string{"message": "Group moved successfully"})
}

// Sync 同步用户组
//
//	@Summary		同步用户组
//	@Description	从外部系统同步用户组
//	@Tags			AuthGroup
//	@Accept			json
//	@Produce		json
//	@Param			body	body		SyncReq	true	"Sync request"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/api/pro/v1/auth/group/sync [post]
//	@Security		bearerAuth
func (h *AuthGroupHandler) Sync(c echo.Context) error {
	var req SyncReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}

	if req.KBId == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	// TODO: 从 LDAP/AD 同步用户组
	h.logger.Info("Sync groups", "kb_id", req.KBId)

	return h.NewResponseWithData(c, map[string]string{"message": "Groups synced successfully"})
}

// Delete 删除用户组
//
//	@Summary		删除用户组
//	@Description	删除指定的用户组
//	@Tags			AuthGroup
//	@Accept			json
//	@Produce		json
//	@Param			id		query		int		true	"Group ID"
//	@Param			kb_id	query		string	true	"Knowledge Base ID"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/api/pro/v1/auth/group/delete [delete]
//	@Security		bearerAuth
func (h *AuthGroupHandler) Delete(c echo.Context) error {
	id := c.QueryParam("id")
	kbID := c.QueryParam("kb_id")

	if id == "" || kbID == "" {
		return h.NewResponseWithError(c, "id and kb_id are required", nil)
	}

	// TODO: 从数据库删除
	h.logger.Info("Delete group", "id", id, "kb_id", kbID)

	return h.NewResponseWithData(c, map[string]string{"message": "Group deleted successfully"})
}

// ============ 数据结构 ============

// GetListResp 用户组列表响应
type GetListResp struct {
	List  []GroupItem `json:"list"`
	Total int64       `json:"total"`
}

// GroupItem 用户组列表项
type GroupItem struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	ParentID  *int    `json:"parent_id,omitempty"`
	Path      string  `json:"path"`
	Position  int     `json:"position"`
	Count     int     `json:"count"`      // 成员数量
	AuthIDs   []int64 `json:"auth_ids"`   // 成员 ID 列表
	CreatedAt string  `json:"created_at"`
}

// GroupTreeNode 用户组树节点
type GroupTreeNode struct {
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	ParentID *int            `json:"parent_id,omitempty"`
	Count    int             `json:"count"`
	Children []GroupTreeNode `json:"children"`
}

// GroupDetail 用户组详情
type GroupDetail struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	ParentID  *int    `json:"parent_id,omitempty"`
	Path      string  `json:"path"`
	Position  int     `json:"position"`
	Members   []int64 `json:"members"`   // 成员 ID 列表
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// CreateReq 创建用户组请求
type CreateReq struct {
	KBId     string  `json:"kb_id" validate:"required"`
	Name     string  `json:"name" validate:"required"`
	ParentID *int    `json:"parent_id,omitempty"`
	AuthIDs  []int64 `json:"auth_ids"` // 成员 ID 列表
}

// CreateResp 创建用户组响应
type CreateResp struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

// UpdateReq 更新用户组请求
type UpdateReq struct {
	ID      int     `json:"id" validate:"required"`
	KBId    string  `json:"kb_id" validate:"required"`
	Name    string  `json:"name"`
	AuthIDs []int64 `json:"auth_ids"` // 成员 ID 列表
}

// MoveReq 移动用户组请求
type MoveReq struct {
	ID       int    `json:"id" validate:"required"`
	KBId     string `json:"kb_id" validate:"required"`
	ParentID *int   `json:"parent_id,omitempty"`
	PrevID   *int   `json:"prev_id,omitempty"`   // 前一个兄弟节点
	NextID   *int   `json:"next_id,omitempty"`   // 后一个兄弟节点
}

// SyncReq 同步用户组请求
type SyncReq struct {
	KBId string `json:"kb_id" validate:"required"`
}
