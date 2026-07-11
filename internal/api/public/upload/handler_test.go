package upload

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KriFinnSher/sany/internal/api/public/upload/mocks"
	entity "github.com/KriFinnSher/sany/internal/entity/upload"
	"github.com/KriFinnSher/sany/internal/logger"
	"go.uber.org/mock/gomock"
)

func TestHandlerServeHTTP(t *testing.T) {
	storeErr := errors.New("storage unavailable")
	tests := []struct {
		name string
		key  string
		file string
		data []byte
		mock func(*mocks.MockUploader)
		code int
		link string
	}{
		{
			name: "uploads file and returns public link",
			key:  "file",
			file: "hello.txt",
			data: []byte("hello, world"),
			mock: func(m *mocks.MockUploader) {
				m.EXPECT().Upload(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, file entity.File) (entity.File, error) {
					if file.Name != "hello.txt" || file.ContentType != "application/octet-stream" || file.Size != int64(len("hello, world")) || string(file.Data) != "hello, world" {
						t.Errorf("unexpected file: %#v", file)
					}
					return entity.File{ID: "file-id"}, nil
				})
			},
			code: http.StatusCreated,
			link: "/api/v1/files?id=file-id",
		},
		{
			name: "rejects missing file form field",
			code: http.StatusBadRequest,
		},
		{
			name: "rejects file larger than limit",
			key:  "file",
			file: "large.bin",
			data: bytes.Repeat([]byte("x"), int(MaxFileSize+(1<<20))),
			code: http.StatusRequestEntityTooLarge,
		},
		{
			name: "returns internal error when upload fails",
			key:  "file",
			file: "hello.txt",
			data: []byte("hello"),
			mock: func(m *mocks.MockUploader) {
				m.EXPECT().Upload(gomock.Any(), gomock.Any()).Return(entity.File{}, storeErr)
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockUploader(ctrl)
			if tt.mock != nil {
				tt.mock(mock)
			}

			w := httptest.NewRecorder()
			New(logger.New(), mock).ServeHTTP(w, multipartRequest(t, tt.key, tt.file, tt.data))

			if w.Code != tt.code {
				t.Fatalf("status = %d, want %d; body = %s", w.Code, tt.code, w.Body.String())
			}
			if tt.link != "" {
				var body response
				if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if body.Link != tt.link {
					t.Errorf("link = %q, want %q", body.Link, tt.link)
				}
			}
		})
	}
}

func multipartRequest(t *testing.T, key, name string, data []byte) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if key != "" {
		part, err := writer.CreateFormFile(key, name)
		if err != nil {
			t.Fatalf("create multipart field: %v", err)
		}
		if _, err := part.Write(data); err != nil {
			t.Fatalf("write multipart field: %v", err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	r := httptest.NewRequest(http.MethodPost, "/api/v1/files", &body)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	return r
}
