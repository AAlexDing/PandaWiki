package pro

import (
	"github.com/google/wire"
)

// ProviderSet is pro handlers providers.
var ProviderSet = wire.NewSet(
	NewContributeHandler,
	NewPromptHandler,
	NewAPITokenHandler,
	NewBlockWordHandler,
	NewAuthHandler,
	NewAuthGroupHandler,
	NewDocumentFeedbackHandler,
	NewCommentModerateHandler,
	NewReleaseHandler,
)
