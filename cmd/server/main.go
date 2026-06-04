package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/KriFinnSher/sany/internal/config"
)

func main() {
	mux := http.NewServeMux()
	cfg := config.MustLoad()

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	ctx := context.Background()
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	log.InfoContext(ctx, "server started", "host", cfg.ServerHost, "port", cfg.ServerPort)

	err := server.ListenAndServe()
	if err != nil {
		log.ErrorContext(ctx, "server stopped", "host", cfg.ServerHost, "port", cfg.ServerPort)
	}
}
