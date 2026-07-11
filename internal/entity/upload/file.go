package upload

import "errors"

var ErrNotFound = errors.New("uploaded file not found")

type File struct {
	ID          string
	Name        string
	ContentType string
	Size        int64
	Data        []byte
}
