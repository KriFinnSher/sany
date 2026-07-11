package upload

import (
	"context"

	entity "github.com/KriFinnSher/sany/internal/entity/upload"
)

type FileUploader interface {
	Upload(context.Context, entity.File) (entity.File, error)
}
