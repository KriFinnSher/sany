package download

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KriFinnSher/sany/internal/api/public/download/mocks"
	entity "github.com/KriFinnSher/sany/internal/entity/upload"
	"github.com/KriFinnSher/sany/internal/logger"
	"go.uber.org/mock/gomock"
)

func TestHandlerServeHTTP(t *testing.T) {
	storeErr := errors.New("storage unavailable")
	tests := []struct {
		name        string
		id          string
		mock        func(*mocks.MockFileGetter)
		code        int
		body        string
		disposition string
	}{
		{
			name: "returns stored ASCII file",
			id:   "file-id",
			mock: func(m *mocks.MockFileGetter) {
				m.EXPECT().Get(gomock.Any(), "file-id").Return(entity.File{
					ID:          "file-id",
					Name:        "hello.txt",
					ContentType: "text/plain",
					Data:        []byte("hello"),
				}, nil)
			},
			code:        http.StatusOK,
			body:        "hello",
			disposition: "inline; filename=hello.txt",
		},
		{
			name: "returns stored Unicode file name",
			id:   "unicode-id",
			mock: func(m *mocks.MockFileGetter) {
				m.EXPECT().Get(gomock.Any(), "unicode-id").Return(entity.File{
					ID:          "unicode-id",
					Name:        "отзыв.docx",
					ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
					Data:        []byte("file data"),
				}, nil)
			},
			code:        http.StatusOK,
			body:        "file data",
			disposition: "inline; filename*=utf-8''%D0%BE%D1%82%D0%B7%D1%8B%D0%B2.docx",
		},
		{
			name: "rejects missing file id",
			code: http.StatusBadRequest,
		},
		{
			name: "returns not found",
			id:   "missing",
			mock: func(m *mocks.MockFileGetter) {
				m.EXPECT().Get(gomock.Any(), "missing").Return(entity.File{}, entity.ErrNotFound)
			},
			code: http.StatusNotFound,
		},
		{
			name: "returns internal error for storage failure",
			id:   "file-id",
			mock: func(m *mocks.MockFileGetter) {
				m.EXPECT().Get(gomock.Any(), "file-id").Return(entity.File{}, storeErr)
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockFileGetter(ctrl)
			if tt.mock != nil {
				tt.mock(mock)
			}

			r := httptest.NewRequest(http.MethodGet, "/api/v1/files?id="+tt.id, nil)
			w := httptest.NewRecorder()
			New(logger.New(), mock).ServeHTTP(w, r)

			if w.Code != tt.code {
				t.Fatalf("status = %d, want %d; body = %s", w.Code, tt.code, w.Body.String())
			}
			if tt.body != "" && w.Body.String() != tt.body {
				t.Errorf("body = %q, want %q", w.Body.String(), tt.body)
			}
			if tt.disposition != "" && w.Header().Get("Content-Disposition") != tt.disposition {
				t.Errorf("Content-Disposition = %q, want %q", w.Header().Get("Content-Disposition"), tt.disposition)
			}
		})
	}
}
