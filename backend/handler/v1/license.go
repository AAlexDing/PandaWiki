package v1

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/consts"
	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/middleware"
)

// LicenseHandler Mock License Handler (用于测试企业版功能)
type LicenseHandler struct {
	*handler.BaseHandler
	logger *log.Logger
}

// NewLicenseHandler 创建 License Handler 并注册路由
func NewLicenseHandler(e *echo.Echo, baseHandler *handler.BaseHandler, logger *log.Logger, auth middleware.AuthMiddleware) *LicenseHandler {
	h := &LicenseHandler{
		BaseHandler: baseHandler,
		logger:      logger.WithModule("handler.v1.license"),
	}

	// 注册路由
	e.GET("/api/v1/license", h.GetLicense, auth.Authorize)

	return h
}

// LicenseResp 授权信息响应
type LicenseResp struct {
	Edition   int32 `json:"edition"`    // 授权版本：0=社区版，1=联创版，2=企业版
	StartedAt int64 `json:"started_at"` // 授权开始时间
	ExpiredAt int64 `json:"expired_at"` // 授权到期时间
	State     int   `json:"state"`      // 授权状态
}

// GetLicense 获取授权信息
//
//	@Summary		Get license
//	@Description	Get license information (Mock for testing)
//	@Tags			license
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	domain.PWResponse{data=LicenseResp}
//	@Router			/api/v1/license [get]
//	@Security		bearerAuth
func (h *LicenseHandler) GetLicense(c echo.Context) error {
	// 返回 Mock 的企业版授权信息
	now := time.Now().Unix()
	resp := LicenseResp{
		Edition:   int32(consts.LicenseEditionEnterprise), // 企业版
		StartedAt: now - 86400*365,                        // 1年前开始
		ExpiredAt: now + 86400*365*10,                     // 10年后到期
		State:     1,                                      // 状态：1=激活
	}

	return h.NewResponseWithData(c, resp)
}
