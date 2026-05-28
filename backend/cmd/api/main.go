package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/findmesh/findmesh/backend/internal/config"
	"github.com/findmesh/findmesh/backend/internal/db"
	"github.com/findmesh/findmesh/backend/internal/server"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	store := db.NewMemoryStore()
	if cfg.Env == "development" {
		db.SeedDemo(store)
	}
	app := server.NewApp(cfg, store, logger)
	logger.Info("findmesh api starting", "addr", cfg.HTTPAddr, "env", cfg.Env)
	if err := http.ListenAndServe(cfg.HTTPAddr, app.Router()); err != nil {
		logger.Error("api stopped", "error", err)
		os.Exit(1)
	}
}
