package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/KriFinnSher/sany/internal/config"
	"github.com/KriFinnSher/sany/internal/logger"
)

func main() {
	mux := http.NewServeMux()
	cfg := config.MustLoad()

	ctx := context.Background()
	logger := logger.New()

	logger.Info(ctx, "server started", "host", cfg.ServerHost, "port", cfg.ServerPort)

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe()
	if err != nil {
		logger.Error(ctx, "server stopped", "err", err, "host", cfg.ServerHost, "port", cfg.ServerPort)
	}
}
