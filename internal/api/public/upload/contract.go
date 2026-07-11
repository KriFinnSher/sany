package upload

import (
	"context"

	entity "github.com/KriFinnSher/sany/internal/entity/upload"
)

type Uploader interface {
	Upload(context.Context, entity.File) (entity.File, error)
}
