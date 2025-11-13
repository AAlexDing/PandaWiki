package v1

import (
	"github.com/google/wire"

	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/middleware"
	"github.com/chaitin/panda-wiki/usecase"
)

type APIHandlers struct {
	UserHandler          *UserHandler
	KnowledgeBaseHandler *KnowledgeBaseHandler
	NodeHandler          *NodeHandler
	AppHandler           *AppHandler
	FileHandler          *FileHandler
	ModelHandler         *ModelHandler
	ConversationHandler  *ConversationHandler
	CrawlerHandler       *CrawlerHandler
	CreationHandler      *CreationHandler
	StatHandler          *StatHandler
	SystemHandler        *SystemHandler
	CommentHandler       *CommentHandler
	AuthV1Handler        *AuthV1Handler
	LicenseHandler       *LicenseHandler
	// Pro handlers 已迁移到 handler/pro 包
	// PromptHandler, BlockWordHandler, APITokenHandler, ContributeHandler 等
	// 现在在 handler/pro 中注册和管理
}

var ProviderSet = wire.NewSet(
	middleware.ProviderSet,
	usecase.ProviderSet,

	handler.NewBaseHandler,
	NewNodeHandler,
	NewAppHandler,
	NewConversationHandler,
	NewUserHandler,
	NewFileHandler,
	NewModelHandler,
	NewKnowledgeBaseHandler,
	NewCrawlerHandler,
	NewCreationHandler,
	NewStatHandler,
	NewSystemHandler,
	NewCommentHandler,
	NewAuthV1Handler,
	NewLicenseHandler,

	wire.Struct(new(APIHandlers), "*"),
)
