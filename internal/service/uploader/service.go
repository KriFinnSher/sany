package uploader

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	entity "github.com/KriFinnSher/sany/internal/entity/upload"
)

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Upload(ctx context.Context, file entity.File) (entity.File, error) {
	if file.Size != int64(len(file.Data)) {
		file.Size = int64(len(file.Data))
	}

	id, err := newID()
	if err != nil {
		return entity.File{}, fmt.Errorf("generate file id: %w", err)
	}
	file.ID = id

	if err := s.repository.Save(ctx, file); err != nil {
		return entity.File{}, fmt.Errorf("save file: %w", err)
	}
	return file, nil
}

func (s *Service) Get(ctx context.Context, id string) (entity.File, error) {
	file, err := s.repository.Get(ctx, id)
	if err != nil {
		return entity.File{}, fmt.Errorf("get file: %w", err)
	}
	return file, nil
}

func newID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
