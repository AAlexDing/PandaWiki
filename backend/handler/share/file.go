package share

import (
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/config"
	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/store/s3"
	"github.com/chaitin/panda-wiki/usecase"
)

type ShareFileHandler struct {
	*handler.BaseHandler
	fileUsecase *usecase.FileUsecase
	minioClient *s3.MinioClient
	config      *config.Config
	logger      *log.Logger
}

func NewShareFileHandler(e *echo.Echo, baseHandler *handler.BaseHandler, fileUsecase *usecase.FileUsecase, minioClient *s3.MinioClient, config *config.Config, logger *log.Logger) *ShareFileHandler {
	h := &ShareFileHandler{
		BaseHandler: baseHandler,
		fileUsecase: fileUsecase,
		minioClient: minioClient,
		config:      config,
		logger:      logger.WithModule("handler.share.file"),
	}

	// 注册路由
	e.POST("/share/pro/v1/file/upload", h.UploadFile)

	return h
}

// UploadFile 上传文件
//
//	@Summary		Upload file
//	@Description	Upload file for contribute
//	@Tags			file
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"File to upload"
//	@Param			kb_id	formData	string	false	"Knowledge Base ID"
//	@Success		200		{object}	domain.PWResponse{data=UploadFileResp}
//	@Router			/share/pro/v1/file/upload [post]
func (h *ShareFileHandler) UploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return h.NewResponseWithError(c, "file is required", err)
	}

	// 获取 kb_id，优先从 form 获取，其次从 header 获取
	kbID := c.FormValue("kb_id")
	if kbID == "" {
		kbID = c.Request().Header.Get("X-KB-ID")
	}
	// 如果都没有，使用 "shared" 作为共享文件目录
	if kbID == "" {
		kbID = "shared"
	}

	// 上传到 MinIO
	key, err := h.fileUsecase.UploadFile(c.Request().Context(), kbID, file)
	if err != nil {
		h.logger.Error("upload file failed", log.Error(err))
		return h.NewResponseWithError(c, "upload file failed", err)
	}

	return h.NewResponseWithData(c, UploadFileResp{
		Key:      key,
		Filename: file.Filename,
		Size:     file.Size,
	})
}

// UploadFileResp 上传文件响应
type UploadFileResp struct {
	Key      string `json:"key"`      // 文件在 MinIO 中的路径，如：kb_id/uuid.ext
	Filename string `json:"filename"` // 原始文件名
	Size     int64  `json:"size"`     // 文件大小
}
