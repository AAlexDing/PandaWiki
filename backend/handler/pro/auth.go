package pro

import (
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/middleware"
)

type AuthHandler struct {
	*handler.BaseHandler
	logger *log.Logger
}

func NewAuthHandler(e *echo.Echo, baseHandler *handler.BaseHandler, logger *log.Logger, auth middleware.AuthMiddleware) *AuthHandler {
	h := &AuthHandler{
		BaseHandler: baseHandler,
		logger:      logger.WithModule("handler.pro.auth"),
	}

	// 注册路由
	e.GET("/api/pro/v1/auth/get", h.GetAuth, auth.Authorize)
	e.POST("/api/pro/v1/auth/set", h.SetAuth, auth.Authorize)
	e.DELETE("/api/pro/v1/auth/delete", h.DeleteAuth, auth.Authorize)

	return h
}

// GetAuth 获取权限
//
//	@Summary		获取权限配置
//	@Description	获取指定知识库的权限配置
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			kb_id	query		string	true	"Knowledge Base ID"
//	@Success		200		{object}	domain.PWResponse{data=GetAuthResp}
//	@Router			/api/pro/v1/auth/get [get]
//	@Security		bearerAuth
func (h *AuthHandler) GetAuth(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	// TODO: 从数据库查询权限配置
	// 目前返回空配置，表示使用默认权限（所有人可读）
	return h.NewResponseWithData(c, GetAuthResp{
		KBId:   kbID,
		Public: true, // 默认公开
		Groups: []AuthGroup{},
		Users:  []AuthUser{},
	})
}

// SetAuth 设置权限
//
//	@Summary		设置权限
//	@Description	设置知识库的权限配置
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		SetAuthReq	true	"Set auth request"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/api/pro/v1/auth/set [post]
//	@Security		bearerAuth
func (h *AuthHandler) SetAuth(c echo.Context) error {
	var req SetAuthReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}

	if req.KBId == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	// TODO: 保存权限配置到数据库
	h.logger.Info("Set auth", "kb_id", req.KBId, "public", req.Public)

	return h.NewResponseWithData(c, map[string]string{"message": "Auth set successfully"})
}

// DeleteAuth 删除权限
//
//	@Summary		删除权限
//	@Description	删除指定用户或用户组的权限
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			kb_id		query		string	true	"Knowledge Base ID"
//	@Param			auth_id		query		int		false	"Auth ID"
//	@Param			group_id	query		int		false	"Group ID"
//	@Success		200			{object}	domain.PWResponse
//	@Router			/api/pro/v1/auth/delete [delete]
//	@Security		bearerAuth
func (h *AuthHandler) DeleteAuth(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	authID := c.QueryParam("auth_id")
	groupID := c.QueryParam("group_id")

	if authID == "" && groupID == "" {
		return h.NewResponseWithError(c, "auth_id or group_id is required", nil)
	}

	// TODO: 从数据库删除权限配置
	h.logger.Info("Delete auth", "kb_id", kbID, "auth_id", authID, "group_id", groupID)

	return h.NewResponseWithData(c, map[string]string{"message": "Auth deleted successfully"})
}

// GetAuthResp 获取权限响应
type GetAuthResp struct {
	KBId   string      `json:"kb_id"`
	Public bool        `json:"public"` // 是否公开
	Groups []AuthGroup `json:"groups"` // 用户组权限
	Users  []AuthUser  `json:"users"`  // 用户权限
}

// AuthGroup 用户组权限
type AuthGroup struct {
	GroupID    int    `json:"group_id"`
	GroupName  string `json:"group_name"`
	Permission string `json:"permission"` // read, write, admin
}

// AuthUser 用户权限
type AuthUser struct {
	AuthID     int    `json:"auth_id"`
	Username   string `json:"username"`
	Permission string `json:"permission"` // read, write, admin
}

// SetAuthReq 设置权限请求
type SetAuthReq struct {
	KBId   string      `json:"kb_id" validate:"required"`
	Public bool        `json:"public"` // 是否公开
	Groups []AuthGroup `json:"groups"` // 用户组权限
	Users  []AuthUser  `json:"users"`  // 用户权限
}
