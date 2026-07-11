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
	fileSaverError := errors.New("file saver unavailable")
	tests := []struct {
		name      string
		configure func(*mocks.MockFileSaver)
		wantErr   bool
	}{
		{
			name: "saves file with generated id and measured size",
			configure: func(mock *mocks.MockFileSaver) {
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
			name: "returns error when file saver fails",
			configure: func(mock *mocks.MockFileSaver) {
				mock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(fileSaverError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			fileSaver := mocks.NewMockFileSaver(ctrl)
			fileGetter := mocks.NewMockFileGetter(ctrl)
			tt.configure(fileSaver)

			file, err := New(fileSaver, fileGetter).Upload(context.Background(), entity.File{Name: "hello.txt", Data: []byte("hello")})
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
	fileGetterError := errors.New("file getter unavailable")
	tests := []struct {
		name      string
		configure func(*mocks.MockFileGetter)
		wantErr   error
	}{
		{
			name: "returns file getter result",
			configure: func(mock *mocks.MockFileGetter) {
				mock.EXPECT().Get(gomock.Any(), "file-id").Return(entity.File{ID: "file-id", Data: []byte("hello")}, nil)
			},
		},
		{
			name: "preserves not found error",
			configure: func(mock *mocks.MockFileGetter) {
				mock.EXPECT().Get(gomock.Any(), "file-id").Return(entity.File{}, entity.ErrNotFound)
			},
			wantErr: entity.ErrNotFound,
		},
		{
			name: "wraps file getter failure",
			configure: func(mock *mocks.MockFileGetter) {
				mock.EXPECT().Get(gomock.Any(), "file-id").Return(entity.File{}, fileGetterError)
			},
			wantErr: fileGetterError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			fileSaver := mocks.NewMockFileSaver(ctrl)
			fileGetter := mocks.NewMockFileGetter(ctrl)
			tt.configure(fileGetter)

			file, err := New(fileSaver, fileGetter).Get(context.Background(), "file-id")
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
