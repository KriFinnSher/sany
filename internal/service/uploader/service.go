package uploader

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	entity "github.com/KriFinnSher/sany/internal/entity/upload"
)

type Service struct {
	fileSaver  FileSaver
	fileGetter FileGetter
}

// New returns a service that stores and retrieves uploaded files.
func New(fileSaver FileSaver, fileGetter FileGetter) *Service {
	return &Service{
		fileSaver:  fileSaver,
		fileGetter: fileGetter,
	}
}

// Upload normalizes file metadata, assigns an ID, and saves the file.
func (s *Service) Upload(ctx context.Context, file entity.File) (entity.File, error) {
	// Data is authoritative because it is persisted with the supplied metadata.
	if file.Size != int64(len(file.Data)) {
		file.Size = int64(len(file.Data))
	}

	id, err := newID()
	if err != nil {
		return entity.File{}, fmt.Errorf("generate file id: %w", err)
	}
	file.ID = id

	if err := s.fileSaver.Save(ctx, file); err != nil {
		return entity.File{}, fmt.Errorf("save file: %w", err)
	}
	return file, nil
}

// Get retrieves a stored file by ID.
func (s *Service) Get(ctx context.Context, id string) (entity.File, error) {
	file, err := s.fileGetter.Get(ctx, id)
	if err != nil {
		return entity.File{}, fmt.Errorf("get file: %w", err)
	}
	return file, nil
}

// newID creates an opaque 128-bit identifier for a stored file.
func newID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
