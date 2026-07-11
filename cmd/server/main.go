package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/KriFinnSher/sany/internal/api/public/download"
	"github.com/KriFinnSher/sany/internal/api/public/upload"
	"github.com/KriFinnSher/sany/internal/config"
	"github.com/KriFinnSher/sany/internal/config/database"
	"github.com/KriFinnSher/sany/internal/logger"
	"github.com/KriFinnSher/sany/internal/service/uploader"
	"github.com/KriFinnSher/sany/internal/storage/sqlite"
)

func main() {
	mux := http.NewServeMux()
	cfg := config.MustLoad()
	db := database.MustLoadSQLite(cfg)

	ctx := context.Background()
	logger := logger.New()

	logger.Info(ctx, "server started", "host", cfg.ServerHost, "port", cfg.ServerPort)

	storage, err := sqlite.New(db)
	if err != nil {
		logger.Error(ctx, "failed to initialize storage", "err", err)
		return
	}
	uploader := uploader.New(storage)

	mux.Handle("POST /api/v1/files", upload.New(logger, uploader))
	mux.Handle("GET /api/v1/files", download.New(logger, uploader))

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort),
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		logger.Error(ctx, "server stopped", "err", err, "host", cfg.ServerHost, "port", cfg.ServerPort)
	}
}
