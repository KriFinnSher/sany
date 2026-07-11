package uploader

import (
	"context"

	entity "github.com/KriFinnSher/sany/internal/entity/upload"
)

type FileSaver interface {
	Save(context.Context, entity.File) error
}

type FileGetter interface {
	Get(ctx context.Context, id string) (entity.File, error)
}
