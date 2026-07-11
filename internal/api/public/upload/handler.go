package upload

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/KriFinnSher/sany/internal/api"
	"github.com/KriFinnSher/sany/internal/logger"
)

const maxFileSize = 50 << 20

type Handler struct {
	log logger.Logger
	// uploader uploader
}

func New(path string, log logger.Logger) *Handler {
	return &Handler{
		log: log.With(logger.OperationField, path),
		// uploader: uploader,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Debug(ctx, "http call")

	if !h.methodPass(r) {
		h.log.Error(ctx, "disallowed method was applied")
	}

	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		h.log.Error(ctx, "failed to parse multipart form", logger.ErrFiled, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(api.NewMessage("bad file received"))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.log.Error(ctx, "failed to get file from form", logger.ErrFiled, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(api.NewMessage("bad form key received"))
		return
	}

	var buf bytes.Buffer

	n, err := io.Copy(&buf, file)
	if err != nil {
		h.log.Error(ctx, "failed to copy file", logger.ErrFiled, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(api.NewMessage("file processing error (copy)"))
		return
	}

	h.log.Debug(ctx, "file was copied from buffer", "bytes", n)

	err = os.WriteFile(header.Filename, buf.Bytes(), 0755)
	if err != nil {
		h.log.Error(ctx, "failed to create file", logger.ErrFiled, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(api.NewMessage("file processing error (create)"))
		return
	}

	h.log.Debug(ctx, "file was created")

	w.Write(api.NewMessage("file created successfully"))
}

func (h *Handler) methodPass(r *http.Request) bool {
	return r.Method == http.MethodPost
}
