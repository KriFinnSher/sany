package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	entity "github.com/KriFinnSher/sany/internal/entity/upload"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) (*Storage, error) {
	storage := &Storage{db: db}
	if err := storage.migrate(context.Background()); err != nil {
		return nil, err
	}
	return storage, nil
}

func (s *Storage) Save(ctx context.Context, file entity.File) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO uploads (id, name, content_type, size, data)
		VALUES (?, ?, ?, ?, ?)`, file.ID, file.Name, file.ContentType, file.Size, file.Data)
	if err != nil {
		return fmt.Errorf("insert upload: %w", err)
	}
	return nil
}

func (s *Storage) Get(ctx context.Context, id string) (entity.File, error) {
	var file entity.File
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, content_type, size, data
		FROM uploads
		WHERE id = ?`, id).Scan(&file.ID, &file.Name, &file.ContentType, &file.Size, &file.Data)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.File{}, entity.ErrNotFound
	}
	if err != nil {
		return entity.File{}, fmt.Errorf("select upload: %w", err)
	}
	return file, nil
}

func (s *Storage) migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS uploads (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			content_type TEXT NOT NULL,
			size INTEGER NOT NULL,
			data BLOB NOT NULL
		)`)
	if err != nil {
		return fmt.Errorf("create uploads table: %w", err)
	}
	return nil
}
