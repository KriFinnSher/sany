package download

import (
	"errors"
	"fmt"
	"mime"
	"net/http"

	"github.com/KriFinnSher/sany/internal/api/http_utils"
	entity "github.com/KriFinnSher/sany/internal/entity/upload"
	"github.com/KriFinnSher/sany/internal/logger"
)

type Handler struct {
	logger     logger.Logger
	fileGetter FileGetter
}

func New(log logger.Logger, fileGetter FileGetter) *Handler {
	return &Handler{
		logger:     log.With(logger.OperationField, "download"),
		fileGetter: fileGetter,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http_utils.WriteError(w, http.StatusBadRequest, "file id is required")
		return
	}

	file, err := h.fileGetter.Get(r.Context(), id)
	if errors.Is(err, entity.ErrNotFound) {
		http_utils.WriteError(w, http.StatusNotFound, "file not found")
		return
	}
	if err != nil {
		h.logger.Error(r.Context(), "get file", logger.ErrFiled, err)
		http_utils.WriteError(w, http.StatusInternalServerError, "failed to get file")
		return
	}

	w.Header().Set("Content-Type", http_utils.ContentType(file.ContentType))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(file.Data)))
	w.Header().Set("Content-Disposition", mime.FormatMediaType("inline", map[string]string{
		"filename": file.Name,
	}))
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(file.Data); err != nil {
		h.logger.Error(r.Context(), "write file", logger.ErrFiled, err)
	}
}
