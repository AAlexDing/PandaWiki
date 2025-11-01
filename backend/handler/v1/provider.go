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
	CommentHandler       *CommentHandler
	AuthV1Handler        *AuthV1Handler
	LicenseHandler       *LicenseHandler
	ProMockHandler       *ProMockHandler
	PromptHandler        *PromptHandler
	BlockWordHandler     *BlockWordHandler
	APITokenHandler      *APITokenHandler
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
	NewCommentHandler,
	NewAuthV1Handler,
	NewLicenseHandler,
	NewProMockHandler,
	NewPromptHandler,
	NewBlockWordHandler,
	NewAPITokenHandler,

	wire.Struct(new(APIHandlers), "*"),
)
