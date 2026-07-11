package download

import (
	"context"

	entity "github.com/KriFinnSher/sany/internal/entity/upload"
)

type FileGetter interface {
	Get(context.Context, string) (entity.File, error)
}
