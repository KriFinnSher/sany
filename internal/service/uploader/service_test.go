package uploader

import (
	"context"
	"errors"
	"testing"

	entity "github.com/KriFinnSher/sany/internal/entity/upload"
	"github.com/KriFinnSher/sany/internal/service/uploader/mocks"
	"go.uber.org/mock/gomock"
)

func TestServiceUpload(t *testing.T) {
	repositoryError := errors.New("repository unavailable")
	tests := []struct {
		name      string
		configure func(*mocks.MockRepository)
		wantErr   bool
	}{
		{
			name: "saves file with generated id and measured size",
			configure: func(mock *mocks.MockRepository) {
				mock.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, file entity.File) error {
					if len(file.ID) != 32 {
						t.Errorf("id length = %d, want 32", len(file.ID))
					}
					if file.Size != 5 {
						t.Errorf("size = %d, want 5", file.Size)
					}
					return nil
				})
			},
		},
		{
			name: "returns error when repository save fails",
			configure: func(mock *mocks.MockRepository) {
				mock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(repositoryError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockRepository(ctrl)
			tt.configure(mock)

			file, err := New(mock).Upload(context.Background(), entity.File{Name: "hello.txt", Data: []byte("hello")})
			if (err != nil) != tt.wantErr {
				t.Fatalf("Upload() error = %v, wantErr %t", err, tt.wantErr)
			}
			if !tt.wantErr && (file.ID == "" || file.Size != 5) {
				t.Errorf("Upload() = %#v, want generated id and size", file)
			}
		})
	}
}

func TestServiceGet(t *testing.T) {
	repositoryError := errors.New("repository unavailable")
	tests := []struct {
		name      string
		configure func(*mocks.MockRepository)
		wantErr   error
	}{
		{
			name: "returns repository file",
			configure: func(mock *mocks.MockRepository) {
				mock.EXPECT().Get(gomock.Any(), "file-id").Return(entity.File{ID: "file-id", Data: []byte("hello")}, nil)
			},
		},
		{
			name: "preserves not found error",
			configure: func(mock *mocks.MockRepository) {
				mock.EXPECT().Get(gomock.Any(), "file-id").Return(entity.File{}, entity.ErrNotFound)
			},
			wantErr: entity.ErrNotFound,
		},
		{
			name: "wraps repository failure",
			configure: func(mock *mocks.MockRepository) {
				mock.EXPECT().Get(gomock.Any(), "file-id").Return(entity.File{}, repositoryError)
			},
			wantErr: repositoryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockRepository(ctrl)
			tt.configure(mock)

			file, err := New(mock).Get(context.Background(), "file-id")
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Get() error = %v, want wrapped %v", err, tt.wantErr)
				}
				return
			}
			if err != nil || file.ID != "file-id" {
				t.Errorf("Get() = %#v, %v; want stored file", file, err)
			}
		})
	}
}
