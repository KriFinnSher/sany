package download

import (
	"context"

	entity "github.com/KriFinnSher/sany/internal/entity/upload"
)

type Getter interface {
	Get(context.Context, string) (entity.File, error)
}
