package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/KriFinnSher/sany/internal/api/public/upload"
	"github.com/KriFinnSher/sany/internal/config"
	"github.com/KriFinnSher/sany/internal/config/database"
	"github.com/KriFinnSher/sany/internal/logger"
	"github.com/KriFinnSher/sany/internal/storage/sqlite"
)

func main() {
	mux := http.NewServeMux()
	cfg := config.MustLoad()
	db := database.MustLoadSQLite(cfg)

	ctx := context.Background()
	logger := logger.New()

	logger.Info(ctx, "server started", "host", cfg.ServerHost, "port", cfg.ServerPort)

	mux.HandleFunc("/api/v1/upload", upload.New("/api/v1/upload", logger).Handle)

	_ = sqlite.New(db)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort),
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		logger.Error(ctx, "server stopped", "err", err, "host", cfg.ServerHost, "port", cfg.ServerPort)
	}
}
