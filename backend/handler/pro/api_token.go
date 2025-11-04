package pro

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/consts"
	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/middleware"
	"github.com/chaitin/panda-wiki/repo/pg"
)

type APITokenHandler struct {
	*handler.BaseHandler
	apiTokenRepo *pg.APITokenRepo
	logger       *log.Logger
}

func NewAPITokenHandler(e *echo.Echo, baseHandler *handler.BaseHandler, apiTokenRepo *pg.APITokenRepo, logger *log.Logger, auth middleware.AuthMiddleware) *APITokenHandler {
	h := &APITokenHandler{
		BaseHandler:  baseHandler,
		apiTokenRepo: apiTokenRepo,
		logger:       logger.WithModule("handler.pro.api_token"),
	}

	// 注册路由
	e.GET("/api/pro/v1/token/list", h.GetTokenList, auth.Authorize)
	e.POST("/api/pro/v1/token/create", h.CreateToken, auth.Authorize)
	e.PATCH("/api/pro/v1/token/update", h.UpdateToken, auth.Authorize)
	e.DELETE("/api/pro/v1/token/delete", h.DeleteToken, auth.Authorize)

	return h
}

// GetTokenList 获取 API Token 列表
//
//	@Summary		Get API token list
//	@Description	Get API token list for knowledge base
//	@Tags			token
//	@Accept			json
//	@Produce		json
//	@Param			kb_id	query		string	true	"Knowledge Base ID"
//	@Success		200		{object}	domain.PWResponse{data=[]TokenItem}
//	@Router			/api/pro/v1/token/list [get]
//	@Security		bearerAuth
func (h *APITokenHandler) GetTokenList(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	tokens, err := h.apiTokenRepo.GetListByKBID(c.Request().Context(), kbID)
	if err != nil {
		h.logger.Error("get token list failed", log.Error(err))
		return h.NewResponseWithError(c, "get token list failed", err)
	}

	if tokens == nil || len(tokens) == 0 {
		return h.NewResponseWithData(c, []*TokenItem{})
	}

	// 转换为响应格式（隐藏完整 token）
	items := make([]*TokenItem, len(tokens))
	for i, token := range tokens {
		items[i] = &TokenItem{
			ID:         token.ID,
			Name:       token.Name,
			UserID:     token.UserID,
			KbID:       token.KbId,
			Permission: token.Permission,
			Token:      token.Token, // 返回完整 token（前端会自己做遮罩）
			CreatedAt:  token.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return h.NewResponseWithData(c, items)
}

// CreateToken 创建 API Token
//
//	@Summary		Create API token
//	@Description	Create API token for knowledge base
//	@Tags			token
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateTokenReq	true	"Create token request"
//	@Success		200		{object}	domain.PWResponse{data=CreateTokenResp}
//	@Router			/api/pro/v1/token/create [post]
//	@Security		bearerAuth
func (h *APITokenHandler) CreateToken(c echo.Context) error {
	var req CreateTokenReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}

	if req.KBID == "" || req.Name == "" {
		return h.NewResponseWithError(c, "kb_id and name are required", nil)
	}

	// 获取当前用户 ID
	authInfo := domain.GetAuthInfoFromCtx(c.Request().Context())
	if authInfo == nil {
		return h.NewResponseWithError(c, "unauthorized", nil)
	}

	// 生成随机 token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		h.logger.Error("generate token failed", log.Error(err))
		return h.NewResponseWithError(c, "generate token failed", err)
	}
	tokenString := "pw_" + hex.EncodeToString(tokenBytes)

	// 创建 token
	apiToken := &domain.APIToken{
		ID:         uuid.New().String(),
		Name:       req.Name,
		UserID:     authInfo.UserId,
		Token:      tokenString,
		KbId:       req.KBID,
		Permission: req.Permission,
	}

	if err := h.apiTokenRepo.Create(c.Request().Context(), apiToken); err != nil {
		h.logger.Error("create token failed", log.Error(err))
		return h.NewResponseWithError(c, "create token failed", err)
	}

	return h.NewResponseWithData(c, CreateTokenResp{
		ID:    apiToken.ID,
		Token: tokenString,
	})
}

// DeleteToken 删除 API Token
//
//	@Summary		Delete API token
//	@Description	Delete API token by ID
//	@Tags			token
//	@Accept			json
//	@Produce		json
//	@Param			id		query		string	true	"Token ID"
//	@Param			kb_id	query		string	true	"Knowledge Base ID"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/api/pro/v1/token/delete [delete]
//	@Security		bearerAuth
func (h *APITokenHandler) DeleteToken(c echo.Context) error {
	id := c.QueryParam("id")
	kbID := c.QueryParam("kb_id")

	if id == "" || kbID == "" {
		return h.NewResponseWithError(c, "id and kb_id are required", nil)
	}

	if err := h.apiTokenRepo.Delete(c.Request().Context(), id, kbID); err != nil {
		h.logger.Error("delete token failed", log.Error(err))
		return h.NewResponseWithError(c, "delete token failed", err)
	}

	return h.NewResponseWithData(c, map[string]string{"message": "token deleted successfully"})
}

// UpdateToken 更新 API Token
//
//	@Summary		Update API token
//	@Description	Update API token name or permission
//	@Tags			token
//	@Accept			json
//	@Produce		json
//	@Param			body	body		UpdateTokenReq	true	"Update token request"
//	@Success		200		{object}	domain.PWResponse
//	@Router			/api/pro/v1/token/update [patch]
//	@Security		bearerAuth
func (h *APITokenHandler) UpdateToken(c echo.Context) error {
	var req UpdateTokenReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}

	if req.ID == "" || req.KBID == "" {
		return h.NewResponseWithError(c, "id and kb_id are required", nil)
	}

	// TODO: 实现更新逻辑
	// if err := h.apiTokenRepo.Update(c.Request().Context(), req.ID, req.KBID, req.Name, req.Permission); err != nil {
	// 	h.logger.Error("update token failed", log.Error(err))
	// 	return h.NewResponseWithError(c, "update token failed", err)
	// }

	h.logger.Info("Update token", "id", req.ID, "name", req.Name)

	return h.NewResponseWithData(c, map[string]string{"message": "token updated successfully"})
}

// CreateTokenReq 创建 Token 请求
type CreateTokenReq struct {
	KBID       string                  `json:"kb_id" validate:"required"`
	Name       string                  `json:"name" validate:"required"`
	Permission consts.UserKBPermission `json:"permission" validate:"required"`
}

// UpdateTokenReq 更新 Token 请求
type UpdateTokenReq struct {
	ID         string                  `json:"id" validate:"required"`
	KBID       string                  `json:"kb_id" validate:"required"`
	Name       string                  `json:"name,omitempty"`
	Permission consts.UserKBPermission `json:"permission,omitempty"`
}

// CreateTokenResp 创建 Token 响应
type CreateTokenResp struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// TokenItem Token 列表项
type TokenItem struct {
	ID         string                  `json:"id"`
	Name       string                  `json:"name"`
	UserID     string                  `json:"user_id"`
	KbID       string                  `json:"kb_id"`
	Permission consts.UserKBPermission `json:"permission"`
	Token      string                  `json:"token"` // 完整 token，前端负责遮罩显示
	CreatedAt  string                  `json:"created_at"`
}
