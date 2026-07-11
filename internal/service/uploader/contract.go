package uploader

import (
	"context"

	entity "github.com/KriFinnSher/sany/internal/entity/upload"
)

type Repository interface {
	Save(ctx context.Context, file entity.File) error
	Get(ctx context.Context, id string) (entity.File, error)
}
