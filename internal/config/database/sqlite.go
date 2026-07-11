package database

import (
	"database/sql"

	"github.com/KriFinnSher/sany/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

func MustLoadSQLite(cfg *config.Config) *sql.DB {
	db, err := sql.Open("sqlite3", cfg.DataSourcePath)
	if err != nil {
		panic(err)
	}

	return db
}
