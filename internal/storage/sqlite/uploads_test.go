package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	entity "github.com/KriFinnSher/sany/internal/entity/upload"
	_ "github.com/mattn/go-sqlite3"
)

func TestStorageSaveAndGet(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	storage, err := New(db)
	if err != nil {
		t.Fatalf("create storage: %v", err)
	}
	file := entity.File{ID: "file-id", Name: "hello.txt", ContentType: "text/plain", Size: 5, Data: []byte("hello")}

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{name: "returns saved file", id: file.ID},
		{name: "returns not found for unknown id", id: "missing", wantErr: entity.ErrNotFound},
	}

	if err := storage.Save(context.Background(), file); err != nil {
		t.Fatalf("save file: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := storage.Get(context.Background(), tt.id)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Get() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if got.ID != file.ID || got.Name != file.Name || got.ContentType != file.ContentType || got.Size != file.Size || string(got.Data) != string(file.Data) {
				t.Errorf("Get() = %#v, want %#v", got, file)
			}
		})
	}
}
