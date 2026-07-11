package upload

import (
	"context"
	"os"
)

type uploader interface {
	Upload(ctx context.Context, file os.File) (link string, err error)
}
