package upload

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"

	"github.com/KriFinnSher/sany/internal/api/http_utils"
	entity "github.com/KriFinnSher/sany/internal/entity/upload"
	"github.com/KriFinnSher/sany/internal/logger"
)

const MaxFileSize int64 = 50 << 20

type Handler struct {
	logger       logger.Logger
	fileUploader FileUploader
}

// New returns a handler for uploading files.
func New(log logger.Logger, fileUploader FileUploader) *Handler {
	return &Handler{
		logger:       log.With(logger.OperationField, "upload"),
		fileUploader: fileUploader,
	}
}

// ServeHTTP bounds the request, validates its multipart file, and stores its contents.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Leave room for multipart overhead while enforcing the file-size limit below.
	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize+(1<<20))
	if err := r.ParseMultipartForm(MaxFileSize); err != nil {
		h.logger.Error(ctx, "parse multipart form", logger.ErrFiled, err)
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			h.logger.Info(ctx, "upload rejected", "reason", "file size limit exceeded")
			http_utils.WriteError(w, http.StatusRequestEntityTooLarge, "file exceeds 50 MiB limit")
			return
		}
		h.logger.Info(ctx, "upload rejected", "reason", "invalid multipart form")
		http_utils.WriteError(w, http.StatusBadRequest, "bad file received")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Error(ctx, "get multipart file", logger.ErrFiled, err)
		h.logger.Info(ctx, "upload rejected", "reason", "missing file form field")
		http_utils.WriteError(w, http.StatusBadRequest, "bad form key received")
		return
	}
	defer file.Close()

	name := filepath.Base(header.Filename)
	contentType := http_utils.ContentType(header.Header.Get("Content-Type"))
	h.logger.Info(ctx, "upload requested", "file_name", name, "content_type", contentType, "size", header.Size)

	if header.Size > MaxFileSize {
		h.logger.Info(ctx, "upload rejected", "file_name", name, "size", header.Size, "reason", "file size limit exceeded")
		http_utils.WriteError(w, http.StatusRequestEntityTooLarge, "file exceeds 50 MiB limit")
		return
	}

	// Read one extra byte so files that exceed the limit are rejected.
	data, err := io.ReadAll(io.LimitReader(file, MaxFileSize+1))
	if err != nil {
		h.logger.Error(ctx, "read uploaded file", logger.ErrFiled, err)
		http_utils.WriteError(w, http.StatusInternalServerError, "file processing error")
		return
	}
	if int64(len(data)) > MaxFileSize {
		h.logger.Info(ctx, "upload rejected", "file_name", name, "size", len(data), "reason", "file size limit exceeded")
		http_utils.WriteError(w, http.StatusRequestEntityTooLarge, "file exceeds 50 MiB limit")
		return
	}

	stored, err := h.fileUploader.Upload(ctx, entity.File{
		// Clients can send paths in multipart names; store only the base filename.
		Name:        name,
		ContentType: contentType,
		Size:        int64(len(data)),
		Data:        data,
	})
	if err != nil {
		h.logger.Error(ctx, "upload file", logger.ErrFiled, err)
		http_utils.WriteError(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	h.logger.Info(ctx, "file uploaded", "file_id", stored.ID, "file_name", stored.Name, "content_type", stored.ContentType, "size", stored.Size)
	http_utils.WriteJSON(w, http.StatusCreated, response{Link: "/api/v1/files?id=" + stored.ID})
}

type response struct {
	Link string `json:"link"`
}
