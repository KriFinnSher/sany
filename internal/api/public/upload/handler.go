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
	log logger.Logger
	up  Uploader
}

func New(log logger.Logger, up Uploader) *Handler {
	return &Handler{
		log: log.With(logger.OperationField, "upload"),
		up:  up,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize+(1<<20))
	if err := r.ParseMultipartForm(MaxFileSize); err != nil {
		h.log.Error(r.Context(), "parse multipart form", logger.ErrFiled, err)
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			http_utils.WriteError(w, http.StatusRequestEntityTooLarge, "file exceeds 50 MiB limit")
			return
		}
		http_utils.WriteError(w, http.StatusBadRequest, "bad file received")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.log.Error(r.Context(), "get multipart file", logger.ErrFiled, err)
		http_utils.WriteError(w, http.StatusBadRequest, "bad form key received")
		return
	}
	defer file.Close()

	if header.Size > MaxFileSize {
		http_utils.WriteError(w, http.StatusRequestEntityTooLarge, "file exceeds 50 MiB limit")
		return
	}

	data, err := io.ReadAll(io.LimitReader(file, MaxFileSize+1))
	if err != nil {
		h.log.Error(r.Context(), "read uploaded file", logger.ErrFiled, err)
		http_utils.WriteError(w, http.StatusInternalServerError, "file processing error")
		return
	}
	if int64(len(data)) > MaxFileSize {
		http_utils.WriteError(w, http.StatusRequestEntityTooLarge, "file exceeds 50 MiB limit")
		return
	}

	stored, err := h.up.Upload(r.Context(), entity.File{
		Name:        filepath.Base(header.Filename),
		ContentType: http_utils.ContentType(header.Header.Get("Content-Type")),
		Size:        int64(len(data)),
		Data:        data,
	})
	if err != nil {
		h.log.Error(r.Context(), "upload file", logger.ErrFiled, err)
		http_utils.WriteError(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	http_utils.WriteJSON(w, http.StatusCreated, response{Link: "/api/v1/files?id=" + stored.ID})
}

type response struct {
	Link string `json:"link"`
}
