package share

import (
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
)

// ShareAuthProHandler Mock Share Auth Pro Handler (用于消除 404 日志)
type ShareAuthProHandler struct {
	*handler.BaseHandler
	logger *log.Logger
}

// NewShareAuthProHandler 创建 Share Auth Pro Handler 并注册路由
func NewShareAuthProHandler(e *echo.Echo, baseHandler *handler.BaseHandler, logger *log.Logger) *ShareAuthProHandler {
	h := &ShareAuthProHandler{
		BaseHandler: baseHandler,
		logger:      logger.WithModule("handler.share.pro.auth"),
	}

	// SSO 认证相关路由
	e.GET("/share/pro/v1/auth/info", h.GetAuthInfo)
	e.POST("/share/pro/v1/auth/cas", h.CASAuth)
	e.POST("/share/pro/v1/auth/dingtalk", h.DingTalkAuth)
	e.POST("/share/pro/v1/auth/feishu", h.FeishuAuth)
	e.POST("/share/pro/v1/auth/github", h.GitHubAuth)
	e.POST("/share/pro/v1/auth/ldap", h.LDAPAuth)
	e.POST("/share/pro/v1/auth/oauth", h.OAuthAuth)
	e.POST("/share/pro/v1/auth/wecom", h.WecomAuth)

	// OAuth 回调路由
	e.GET("/share/pro/v1/openapi/cas/callback", h.CASCallback)
	e.GET("/share/pro/v1/openapi/dingtalk/callback", h.DingTalkCallback)
	e.GET("/share/pro/v1/openapi/feishu/callback", h.FeishuCallback)
	e.GET("/share/pro/v1/openapi/github/callback", h.GitHubCallback)
	e.GET("/share/pro/v1/openapi/oauth/callback", h.OAuthCallback)
	e.GET("/share/pro/v1/openapi/wecom/callback", h.WecomCallback)

	// 文档反馈
	e.POST("/share/pro/v1/document/feedback", h.DocumentFeedback)

	return h
}

// GetAuthInfo 获取认证信息（空实现，消除 404）
func (h *ShareAuthProHandler) GetAuthInfo(c echo.Context) error {
	// 返回空对象，前端会理解为"没有高级认证配置"
	return h.NewResponseWithData(c, map[string]interface{}{})
}

// ============ SSO 认证方法 ============

// CASAuth CAS 认证
//
//	@Summary		CAS 认证
//	@Tags			Auth
//	@Router			/share/pro/v1/auth/cas [post]
func (h *ShareAuthProHandler) CASAuth(c echo.Context) error {
	return h.NewResponseWithError(c, "Pro feature not available", nil)
}

// DingTalkAuth 钉钉认证
//
//	@Summary		钉钉认证
//	@Tags			Auth
//	@Router			/share/pro/v1/auth/dingtalk [post]
func (h *ShareAuthProHandler) DingTalkAuth(c echo.Context) error {
	return h.NewResponseWithError(c, "Pro feature not available", nil)
}

// FeishuAuth 飞书认证
//
//	@Summary		飞书认证
//	@Tags			Auth
//	@Router			/share/pro/v1/auth/feishu [post]
func (h *ShareAuthProHandler) FeishuAuth(c echo.Context) error {
	return h.NewResponseWithError(c, "Pro feature not available", nil)
}

// GitHubAuth GitHub 认证
//
//	@Summary		GitHub 认证
//	@Tags			Auth
//	@Router			/share/pro/v1/auth/github [post]
func (h *ShareAuthProHandler) GitHubAuth(c echo.Context) error {
	return h.NewResponseWithError(c, "Pro feature not available", nil)
}

// LDAPAuth LDAP 认证
//
//	@Summary		LDAP 认证
//	@Tags			Auth
//	@Router			/share/pro/v1/auth/ldap [post]
func (h *ShareAuthProHandler) LDAPAuth(c echo.Context) error {
	return h.NewResponseWithError(c, "Pro feature not available", nil)
}

// OAuthAuth OAuth 认证
//
//	@Summary		OAuth 认证
//	@Tags			Auth
//	@Router			/share/pro/v1/auth/oauth [post]
func (h *ShareAuthProHandler) OAuthAuth(c echo.Context) error {
	return h.NewResponseWithError(c, "Pro feature not available", nil)
}

// WecomAuth 企业微信认证
//
//	@Summary		企业微信认证
//	@Tags			Auth
//	@Router			/share/pro/v1/auth/wecom [post]
func (h *ShareAuthProHandler) WecomAuth(c echo.Context) error {
	return h.NewResponseWithError(c, "Pro feature not available", nil)
}

// ============ OAuth 回调方法 ============

// CASCallback CAS 回调
//
//	@Summary		CAS OAuth 回调
//	@Tags			Auth
//	@Router			/share/pro/v1/openapi/cas/callback [get]
func (h *ShareAuthProHandler) CASCallback(c echo.Context) error {
	return c.String(200, "Pro feature not available")
}

// DingTalkCallback 钉钉回调
//
//	@Summary		钉钉 OAuth 回调
//	@Tags			Auth
//	@Router			/share/pro/v1/openapi/dingtalk/callback [get]
func (h *ShareAuthProHandler) DingTalkCallback(c echo.Context) error {
	return c.String(200, "Pro feature not available")
}

// FeishuCallback 飞书回调
//
//	@Summary		飞书 OAuth 回调
//	@Tags			Auth
//	@Router			/share/pro/v1/openapi/feishu/callback [get]
func (h *ShareAuthProHandler) FeishuCallback(c echo.Context) error {
	return c.String(200, "Pro feature not available")
}

// GitHubCallback GitHub 回调
//
//	@Summary		GitHub OAuth 回调
//	@Tags			Auth
//	@Router			/share/pro/v1/openapi/github/callback [get]
func (h *ShareAuthProHandler) GitHubCallback(c echo.Context) error {
	return c.String(200, "Pro feature not available")
}

// OAuthCallback OAuth 回调
//
//	@Summary		OAuth 回调
//	@Tags			Auth
//	@Router			/share/pro/v1/openapi/oauth/callback [get]
func (h *ShareAuthProHandler) OAuthCallback(c echo.Context) error {
	return c.String(200, "Pro feature not available")
}

// WecomCallback 企业微信回调
//
//	@Summary		企业微信 OAuth 回调
//	@Tags			Auth
//	@Router			/share/pro/v1/openapi/wecom/callback [get]
func (h *ShareAuthProHandler) WecomCallback(c echo.Context) error {
	return c.String(200, "Pro feature not available")
}

// ============ 其他功能 ============

// DocumentFeedback 文档反馈
//
//	@Summary		文档反馈
//	@Tags			Document
//	@Router			/share/pro/v1/document/feedback [post]
func (h *ShareAuthProHandler) DocumentFeedback(c echo.Context) error {
	return h.NewResponseWithData(c, map[string]string{"message": "Feedback received"})
}
